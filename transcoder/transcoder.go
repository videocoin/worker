package transcoder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/dispatcher/v1"
	syncerv1 "github.com/videocoin/cloud-api/syncer/v1"
	validatorv1 "github.com/videocoin/cloud-api/validator/v1"
	"github.com/videocoin/cloud-pkg/retry"
	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/stream"
	"github.com/videocoin/transcode/transcoder/hlswatcher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Transcoder struct {
	logger     *logrus.Entry
	t          *time.Ticker
	clientID   string
	dispatcher v1.DispatcherServiceClient
	outputDir  string
	cmd        *exec.Cmd
	caller     *caller.Caller
	sc         *stream.StreamClient
	task       *v1.Task
	syncerAddr string
	HLSWatcher *hlswatcher.Watcher
}

func NewTranscoder(
	logger *logrus.Entry,
	dispatcher v1.DispatcherServiceClient,
	clientID string,
	outputDir string,
	caller *caller.Caller,
	syncerAddr string,
) (*Transcoder, error) {
	return &Transcoder{
		logger:     logger,
		clientID:   clientID,
		dispatcher: dispatcher,
		outputDir:  outputDir,
		caller:     caller,
		sc:         &stream.StreamClient{},
		syncerAddr: syncerAddr,
		HLSWatcher: hlswatcher.New(time.Second * 2),
	}, nil
}

func (t *Transcoder) Start() error {
	t.logger.Infof("starting transcoder")
	t.t = time.NewTicker(5 * time.Second)
	t.dispatch()
	return nil
}

func (t *Transcoder) Stop() error {
	t.logger.Infof("stopping transcoder")
	t.t.Stop()
	return nil
}

func (t *Transcoder) dispatch() error {
	req := &v1.TaskPendingRequest{ClientID: t.clientID}

	for range t.t.C {
		t.logger.Infof("waiting task...")

		ctx := context.Background()
		task, err := t.dispatcher.GetPendingTask(ctx, req)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				if st.Code() == codes.NotFound {
					t.logger.Infof("no task")
					continue
				}
			}
			t.logger.Errorf("failed to get task: %s", err)
		}

		if task == nil || task.ID == "" {
			t.logger.Infof("no task")
			continue
		}

		t.task = task

		err = t.runTask(task)
		if err != nil {
			t.logger.
				WithField("task_id", task.ID).
				Errorf("failed to transcode: %s", err)

			t.dispatcher.MarkTaskAsFailed(context.Background(), &v1.TaskRequest{
				ClientID: t.clientID,
				ID:       t.task.ID,
			})

			t.task = nil

			continue
		}

		t.task = nil
	}

	return nil
}

func (t *Transcoder) runTask(task *v1.Task) error {
	logger := t.logger.WithField("task_id", task.ID)
	logger.Debugf("task: %+v", task)

	logger.Info("running task")

	task.Cmdline = strings.Replace(task.Cmdline, "$OUTPUT", t.outputDir, -1)
	task.Output.Path = strings.Replace(task.Output.Path, "$OUTPUT", t.outputDir, -1)

	if task.ProfileID == "test" {
		task.Input.URI = "/tmp/in.mp4"
	}

	if _, err := os.Stat(task.Output.Path); os.IsNotExist(err) {
		logger.Debugf("creating dir: %s", task.Output.Path)

		os.RemoveAll(task.Output.Path)
		mkdirErr := os.Mkdir(task.Output.Path, 0777)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	err := retry.RetryWithAttempts(5, time.Second*10, func() error {
		logger.Infof("checking source %s", task.Input.URI)
		return checkSource(task.Input.URI)
	})
	if err != nil {
		return err
	}

	c := strings.Split(task.Cmdline, " ")
	cmdName := c[0]
	cmdArgs := c[1:]
	t.cmd = exec.Command(cmdName, cmdArgs...)

	wg := &sync.WaitGroup{}
	ffmpegErrCh := make(chan error, 1)

	taskStatCtx, taskStatCancel := context.WithCancel(context.Background())
	hlssrCtx, hlssrCancel := context.WithCancel(context.Background())

	streamClient, err := stream.NewStreamClient(
		task.StreamContractAddress,
		t.caller)
	if err != nil {
		return err
	}

	t.sc = streamClient

	wg.Add(1)
	go t.runFFmpeg(task, wg, ffmpegErrCh)

	wg.Add(1)
	go t.runTaskStatMonitor(taskStatCtx, task, wg, taskStatCancel)

	t.hlsFlow(
		taskStatCtx,
		taskStatCancel,
		hlssrCtx,
		hlssrCancel,
		wg,
		task,
		t.OnSegmentTranscoded,
	)

	err = <-ffmpegErrCh
	if err != nil {
		taskStatCancel()
		hlssrCancel()

		t.dispatcher.MarkTaskAsFailed(context.Background(), &v1.TaskRequest{
			ClientID: t.clientID,
			ID:       t.task.ID,
		})

		return err
	}

	taskStatCancel()

	wg.Wait()

	t.dispatcher.MarkTaskAsCompleted(context.Background(), &v1.TaskRequest{
		ClientID: t.clientID,
		ID:       t.task.ID,
	})

	logger.Info("task has been completed")

	return nil
}

func (t *Transcoder) OnSegmentTranscoded(segment *hlswatcher.SegmentInfo) error {
	logger := t.logger.WithFields(logrus.Fields{
		"task_id": t.task.ID,
		"segment": segment.Num,
	})

	logger.Info("segment has been transcoded")

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Upload segment
	logger.Info("uploading segment")

	err := retry.RetryWithAttempts(5, time.Second*1, func() error {
		return t.uploadSegmentViaHttp(t.task, segment)
	})
	if err != nil {
		logger.Errorf("failed to upload segment: %s", err)
		return err
	}

	logger.Info("segment has been uploaded")

	//

	idx := -1
	counter := 0
	for {
		chunks, err := t.sc.GetInChunks()
		if err != nil {
			logger.Errorf("failed to get in chunks: %s", err)
			return err
		}

		logger.Debugf("GetInChunks: %+v\n", chunks)

		if len(chunks) > 0 {
			idx = SearchBigInts(chunks, big.NewInt(int64(segment.Num)))
			if idx >= 0 {
				break
			}
			if counter >= 30 {
				err = fmt.Errorf("failed to search in chunks: %s", err)
				logger.Error(err)
				return err
			}
			counter++
		}

		time.Sleep(time.Second * 2)
	}

	if idx >= 0 {
		logger.Info("submitting proof")

		inChunkID := big.NewInt(int64(segment.Num))
		outChunkID := inChunkID

		profileID := new(big.Int)
		profiles, _ := t.sc.GetProfiles()
		if len(profiles) > 0 {
			profileID = profiles[0]
		}

		tx, err := t.sc.SubmitProof(inChunkID, outChunkID, profileID)
		if err != nil {
			logger.Errorf("failed to submit proof: %s", err)
			return err
		}

		logger.Debugf("submitting proof tx %+v\n", tx.Hash().String())

		if t.task != nil {
			logger.Info("validating proof")

			ctx := context.Background()
			vpReq := &validatorv1.ValidateProofRequest{
				StreamContractAddress: t.task.StreamContractAddress,
				ProfileId:             profileID.Bytes(),
				OutputChunkId:         outChunkID.Bytes(),
				StreamId:              t.task.ID,
			}
			_, err = t.dispatcher.ValidateProof(ctx, vpReq)
			if err != nil {
				logger.Errorf("failed to validate proof: %s", err)
				return err
			}
		}
	}

	return nil
}

func (t *Transcoder) uploadSegment(task *v1.Task, segment *hlswatcher.SegmentInfo) error {
	t.logger.
		WithField("segment", segment.Num).
		WithField("path", segment.Source).
		Debug("uploading segment")

	f, err := os.Open(segment.Source)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	req := &syncerv1.SyncRequest{
		Path:        fmt.Sprintf("%s/%d.ts", task.ID, segment.Num),
		ContentType: "video/MP2T",
		Data:        data,
		Duration:    segment.Duration,
	}

	dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer dcancel()

	_, err = t.dispatcher.Sync(dctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transcoder) uploadSegmentViaHttp(task *v1.Task, segment *hlswatcher.SegmentInfo) error {
	t.logger.
		WithField("segment", segment.Num).
		WithField("path", segment.Source).
		Debug("uploading segment via http")

	params := map[string]string{
		"path":        fmt.Sprintf("%s/%d.ts", task.ID, segment.Num),
		"ct":          "video/MP2T",
		"segment_num": strconv.FormatInt(int64(segment.Num), 10),
	}

	request, err := newfileUploadRequest(t.syncerAddr, params, "file", segment.Source)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return nil
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
