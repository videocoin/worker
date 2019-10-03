package transcoder

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/dispatcher/v1"
	"github.com/videocoin/cloud-pkg/retry"
)

func (t *Transcoder) hlsFlow(
	jobStatCtx context.Context,
	jobStatCancel context.CancelFunc,
	hlssrCtx context.Context,
	hlssrCancel context.CancelFunc,
	wg *sync.WaitGroup,
	task *v1.Task,
) {
	t.logger.Debug("starting hls watcher")

	time.Sleep(time.Second * 2)

	go t.HLSWatcher.Start()
	t.HLSWatcher.Wait()
	t.watchAllHLSOutput(task)

	t.logger.Debugf("watch playlists: %+v", t.HLSWatcher.Files())

	time.Sleep(time.Second * 2)

	wg.Add(1)
	go t.runHLSSegmentReciever(hlssrCtx, task, wg, hlssrCancel)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-jobStatCtx.Done():
				t.logger.Debug("stopping hls watcher")
				t.HLSWatcher.Stop()
				hlssrCancel()
			case <-hlssrCtx.Done():
				t.logger.Debug("stopping hls watcher")
				t.HLSWatcher.Stop()
				jobStatCancel()
				return
			}
		}
	}()
}

func (t *Transcoder) watchAllHLSOutput(task *v1.Task) {
	path := task.Output.Path + "/index.m3u8"
	err := t.HLSWatcher.Add(path)
	if err != nil {
		t.logger.Errorf(
			"failed to add hls playlist %s to watcher: %s",
			path,
			err,
		)
	}
}

func (t *Transcoder) unwatchAllHLSOutput(task *v1.Task) {
	t.HLSWatcher.Remove(task.Output.Path + "/index.m3u8")
}

func (t *Transcoder) runHLSSegmentReciever(
	ctx context.Context,
	task *v1.Task,
	wg *sync.WaitGroup,
	cancel context.CancelFunc,
) {
	defer func() {
		t.logger.Debug("hls segment reciever has been stopped")
		cancel()
		wg.Done()
	}()

	t.logger.Debug("starting hls segment reciever")

	variantsCount := len(t.HLSWatcher.Files())
	lastSegmentCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case segment := <-t.HLSWatcher.SegmentsCh:
			if segment == nil {
				continue
			}

			t.logger.WithFields(logrus.Fields{
				"segment": segment.Num,
				"source":  segment.Source,
			}).Debug("segment has been transcoded")

			err := retry.RetryWithAttempts(5, time.Millisecond*200, func() error {
				// return t.uploadSegment(task, segment)
				return nil
			})
			if err != nil {
				t.logger.
					WithFields(logrus.Fields{
						"segment": segment.Num,
						"source":  segment.Source,
					}).
					Errorf("failed to upload segment: %s", err)
				return
			}

			if segment.IsLast {
				lastSegmentCount++

				t.logger.Debugf("variants: %d/%d", lastSegmentCount, variantsCount)

				if lastSegmentCount == variantsCount {
					t.logger.Debug("all playlist variants have been processed")
					return
				}
			}
		case err := <-t.HLSWatcher.ErrCh:
			t.logger.Error(err)
		}
	}
}
