package transcoder

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	"github.com/videocoin/cloud-pkg/retry"
	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/stream"
	"github.com/videocoin/transcode/transcoder/hlswatcher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Transcoder struct {
	logger         *logrus.Entry
	t              *time.Ticker
	clientID       string
	dispatcher     v1.DispatcherServiceClient
	outputDir      string
	cmd            *exec.Cmd
	caller         *caller.Caller
	sc             *stream.Client
	task           *v1.Task
	lastSegmentNum uint64
	syncerAddr     string
	HLSWatcher     *hlswatcher.Watcher
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
		sc:         &stream.Client{},
		syncerAddr: syncerAddr,
		HLSWatcher: hlswatcher.New(time.Second * 2),
	}, nil
}

func (t *Transcoder) Start() error {
	t.logger.Infof("starting transcoder")
	t.t = time.NewTicker(5 * time.Second)
	err := t.dispatch()
	return err
}

func (t *Transcoder) Stop() error {
	t.logger.Infof("stopping transcoder")
	if t.t != nil {
		t.t.Stop()
	}
	return nil
}

func (t *Transcoder) IsRunning() bool {
	return t.t != nil
}

func (t *Transcoder) IsWorking() bool {
	return t.task != nil
}

func (t *Transcoder) dispatch() error {
	req := &v1.TaskPendingRequest{ClientID: t.clientID}

	for range t.t.C {
		t.logger.Info("waiting task...")

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

		err = t.runTask()
		if err != nil {
			t.logger.
				WithField("task_id", task.ID).
				Errorf("failed to transcode: %s", err)

			t.logger.Error(t.dispatcher.MarkTaskAsFailed(context.Background(), &v1.TaskRequest{
				ClientID: t.clientID,
				ID:       t.task.ID,
			}))

			t.task = nil

			continue
		}

		_, err = t.dispatcher.MarkTaskAsCompleted(context.Background(), &v1.TaskRequest{
			ClientID: t.clientID,
			ID:       t.task.ID,
		})
		if err != nil {
			t.logger.
				WithField("task_id", task.ID).
				Errorf("failed to complete: %s", err)

		}

		t.task = nil
	}

	return nil
}

func (t *Transcoder) preRunTask() error {
	logger := t.logger.WithField("task_id", t.task.ID)

	t.task.Cmdline = strings.Join(strings.Fields(strings.TrimSpace(t.task.Cmdline)), " ")
	t.task.Cmdline = strings.Replace(t.task.Cmdline, "$OUTPUT", t.outputDir, -1)
	t.task.Output.Path = strings.Replace(t.task.Output.Path, "$OUTPUT", t.outputDir, -1)

	if _, err := os.Stat(t.task.Output.Path); os.IsNotExist(err) {
		logger.Debugf("creating dir: %s", t.task.Output.Path)

		os.RemoveAll(t.task.Output.Path)
		mkdirErr := os.Mkdir(t.task.Output.Path, 0777)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	err := retry.RetryWithAttempts(5, time.Second*10, func() error {
		logger.Infof("checking source %s", t.task.Input.URI)
		return checkSource(t.task.Input.URI)
	})
	if err != nil {
		return err
	}

	c := strings.Split(t.task.Cmdline, " ")
	cmdName := c[0]
	cmdArgs := c[1:]
	t.cmd = exec.Command(cmdName, cmdArgs...)

	sc, err := stream.NewClient(t.task.StreamContractAddress, t.caller)
	if err != nil {
		return err
	}

	t.sc = sc

	return nil
}

func (t *Transcoder) runTask() error {
	logger := t.logger.WithField("task_id", t.task.ID)
	logger.Debugf("task: %+v", t.task)
	logger.Info("running task")

	err := t.preRunTask()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)

	tmCtx, tmCancel := context.WithCancel(context.Background())
	hlsCtx, hlsCancel := context.WithCancel(context.Background())

	wg.Add(1)
	go t.runFFmpeg(wg, errCh)

	wg.Add(1)
	go t.runTaskMonitor(tmCtx, wg, tmCancel)

	if t.task.IsOutputHLS() {
		t.hlsFlow(tmCtx, tmCancel, hlsCtx, hlsCancel, wg, t.OnSegmentTranscoded)
	}

	err = <-errCh
	if err != nil {
		tmCancel()
		hlsCancel()
		return err
	}

	tmCancel()

	wg.Wait()

	if t.task.IsOutputHLS() && t.lastSegmentNum > 0 {
		logger.Debug("uploading segment file")

		chunks, err := t.sc.GetInChunks()

		logger.Debugf("chunks: %+v", chunks)
		logger.Debugf("chunks err: %+v", err)

		if err == nil {
			lastChunkNum := chunks[len(chunks)-1]

			logger.Debugf("last chunk num: %+v", lastChunkNum.Uint64())
			logger.Debugf("last processed segment: %+v", t.lastSegmentNum)

			outputPath := t.task.Output.Path + "/index.m3u8"
			segments, _ := t.HLSWatcher.ExtractSegments(outputPath)

			for _, s := range segments {
				logger.Debugf("segment: %+v", s)
			}

			if len(segments) > 0 {
				for _, segment := range segments {
					if segment.Num > t.lastSegmentNum && segment.Num <= lastChunkNum.Uint64() {

						logger.
							WithField("segment", segment.Num).
							Debug("uploading last segments")

						if segment.Num == lastChunkNum.Uint64() {
							segment.IsLast = true
						}

						err := t.OnSegmentTranscoded(segment)
						if err != nil {
							logger.
								WithField("segment", segment.Num).
								Debugf("failed to call OnSegmentTranscoded for last segments: %s", err)
						}
					}
				}
			}
		}
	}

	if t.task.IsOutputFile() {
		logger.Debug("uploading single segment file")

		segment := &hlswatcher.SegmentInfo{
			Source:   t.task.Output.Path + "/" + t.task.Output.Name,
			Num:      uint64(t.task.Output.Num),
			Name:     t.task.Output.Name,
			Duration: t.task.Output.Duration,
			IsVOD:    true,
		}
		err := t.OnSegmentTranscoded(segment)
		if err != nil {
			logger.Error(err)
		}
	}
	tmCancel()
	hlsCancel()
	logger.Info("task has been completed")

	return nil
}

func (t *Transcoder) OnSegmentTranscoded(segment *hlswatcher.SegmentInfo) error {
	logger := t.logger.WithFields(logrus.Fields{
		"task_id": t.task.ID,
		"segment": segment.Num,
	})

	logger.Info("segment has been transcoded")

	segmentReq := &v1.SegmentRequest{
		TaskID:    t.task.ID,
		StreamID:  t.task.StreamID,
		ClientID:  t.task.ClientID,
		ProfileID: t.task.ProfileID,
		UserID:    t.task.OwnerID,
		Num:       segment.Num,
		Duration:  segment.Duration,
	}
	_, err := t.dispatcher.MarkSegmentAsTranscoded(context.Background(), segmentReq)
	if err != nil {
		logger.Debugf("[ERR] failed to mark segment as transcoded: %s", err)
	}

	logger.Debug("uploading segment")

	err = t.uploadSegment(segment)
	if err != nil {
		return err
	}

	t.lastSegmentNum = segment.Num

	ok, err := t.waitGetInChunks(segment.Num)
	if err != nil {
		return err
	}

	if ok {
		logger.Debug("validate proof")
		err := t.submitAndValidateProof(segment)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Transcoder) uploadSegmentViaHTTP(task *v1.Task, segment *hlswatcher.SegmentInfo) error {
	t.logger.
		WithField("segment", segment.Num).
		WithField("path", segment.Source).
		Debug("uploading segment via http")

	params := map[string]string{
		"path":        fmt.Sprintf("%s/%d.ts", task.StreamID, segment.Num),
		"ct":          "video/MP2T",
		"segment_num": strconv.FormatInt(int64(segment.Num), 10),
		"duration":    fmt.Sprintf("%f", segment.Duration),
	}

	if segment.IsLast {
		params["last"] = "y"
	}

	if segment.IsVOD {
		params["vod"] = "y"
	}

	defer func() {
		if err := os.Remove(segment.Source); err != nil {
			t.logger.Errorf("failed to delete segment: %s", err)
		}
	}()

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
	if err != nil {
		return nil, err
	}
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
