package transcoder

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/dispatcher/v1"
	"github.com/videocoin/cloud-pkg/retry"
	"github.com/videocoin/transcode/transcoder/hlswatcher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Transcoder struct {
	logger     *logrus.Entry
	t          *time.Ticker
	machineID  string
	dispatcher v1.DispatcherServiceClient
	outputDir  string
	cmd        *exec.Cmd
	HLSWatcher *hlswatcher.Watcher
}

func NewTranscoder(
	logger *logrus.Entry,
	dispatcher v1.DispatcherServiceClient,
	machineID string,
	outputDir string,
) (*Transcoder, error) {
	return &Transcoder{
		logger:     logger,
		machineID:  machineID,
		dispatcher: dispatcher,
		outputDir:  outputDir,
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
	req := &v1.TaskPendingRequest{MachineID: t.machineID}

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

		err = t.runTask(task)
		if err != nil {
			t.logger.
				WithField("task_id", task.ID).
				Errorf("failed to transcode: %s", err)
			continue
		}
	}

	return nil
}

func (t *Transcoder) runTask(task *v1.Task) error {
	logger := t.logger.WithField("id", task.ID)
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

	err := retry.RetryWithAttempts(30, time.Second*1, func() error {
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

	wg.Add(1)
	go t.runFFmpeg(task, wg, ffmpegErrCh)

	wg.Add(1)
	go t.runTaskStatMonitor(taskStatCtx, task, wg, taskStatCancel)

	t.hlsFlow(taskStatCtx, taskStatCancel, hlssrCtx, hlssrCancel, wg, task)

	err = <-ffmpegErrCh
	if err != nil {
		return err
	}

	wg.Wait()

	logger.Info("task has been completed")

	return nil
}
