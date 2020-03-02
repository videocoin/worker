package transcoder

import (
	"context"
	"sync"
	"time"

	v1 "github.com/videocoin/cloud-api/dispatcher/v1"
	"github.com/videocoin/transcode/transcoder/hlswatcher"
)

type SegmentRecvFunc func(*hlswatcher.SegmentInfo) error

func (t *Transcoder) hlsFlow(  //nolint
	jobStatCtx context.Context,
	jobStatCancel context.CancelFunc,
	hlssrCtx context.Context,
	hlssrCancel context.CancelFunc,
	wg *sync.WaitGroup,
	callback SegmentRecvFunc,
) {
	t.logger.Debug("starting hls watcher")

	time.Sleep(time.Second * 2)

	go t.logger.Error(t.HLSWatcher.Start())
	t.HLSWatcher.Wait()
	t.watchAllHLSOutput(t.task)

	t.logger.Debugf("watch playlists: %+v", t.HLSWatcher.Files())

	time.Sleep(time.Second * 2)

	wg.Add(1)
	go t.runHLSSegmentReceiver(hlssrCtx, wg, hlssrCancel, callback)

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

func (t *Transcoder) runHLSSegmentReceiver(
	ctx context.Context,
	wg *sync.WaitGroup,
	cancel context.CancelFunc,
	callback SegmentRecvFunc,
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

			err := callback(segment)
			if err != nil {
				t.logger.Error(t.dispatcher.MarkTaskAsFailed(context.Background(), &v1.TaskRequest{
					ClientID: t.clientID,
					ID:       t.task.ID,
				}))
				return
			}

			if segment.IsLast {
				lastSegmentCount++
				if lastSegmentCount == variantsCount {
					return
				}
			}
		case err := <-t.HLSWatcher.ErrCh:
			t.logger.Error(err)
		}
	}
}
