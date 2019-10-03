package transcoder

import (
	"context"
	"os"
	"sync"
	"time"

	v1 "github.com/videocoin/cloud-api/dispatcher/v1"
)

func (t *Transcoder) runTaskStatMonitor(
	ctx context.Context,
	task *v1.Task,
	wg *sync.WaitGroup,
	cancel context.CancelFunc,
) {
	defer func() {
		t.logger.Debug("task stat monitor has been stopped")
		cancel()
		wg.Done()
	}()

	t.logger.Debug("starting task stat monitor")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ctx := context.Background()
			taskResp, err := t.dispatcher.GetTask(
				ctx,
				&v1.TaskRequest{ID: task.ID},
			)
			if err != nil {
				t.logger.Debugf("[WARN] failed to get task: %s", err)
				continue
			}

			t.logger.Debugf("task status is %s", task.Status.String())

			switch taskResp.Status {
			case v1.TaskStatusCanceled:
				{
					if t.cmd == nil || t.cmd.Process == nil {
						t.logger.Debug("[WARN] ffmpeg already terminated")
						return
					}

					perr := t.cmd.Process.Signal(os.Interrupt)
					if perr != nil {
						if perr == ErrProcessAlreadyFinished {
							t.logger.Debugf("[WARN] ffmpeg: failed to interrupt: %s", perr)
							return
						}

						t.logger.Errorf("ffmpeg: failed to interrupt: %s", perr)
						return
					}

					t.logger.Debug("ffmpeg has been terminated")
					return
				}
			case v1.TaskStatusFailed:
				{
					if t.cmd == nil || t.cmd.Process == nil {
						t.logger.Debugf("[WARN] ffmpeg already terminated")
						return
					}

					perr := t.cmd.Process.Kill()
					if perr != nil {
						if perr == ErrProcessAlreadyFinished {
							t.logger.Warningf("ffmpeg: failed to kill: %s", perr)
							return
						}

						t.logger.Errorf("ffmpeg: failed to kill: %s", perr)
						return
					}

					t.logger.Info("ffmpeg has been killed")
					return
				}
			case v1.TaskStatusCompleted:
				return
			}
		}
	}
}
