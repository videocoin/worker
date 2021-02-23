package transcoder

import (
	"bytes"
	"fmt"
	"sync"
)

func (t *Transcoder) runFFmpeg(wg *sync.WaitGroup, errCh chan error) {
	stopCh := make(chan bool, 1)
	defer func() {
		close(errCh)
		wg.Done()
		stopCh <- true
		close(stopCh)
		t.logger.Debug("ffmpeg has been completed")
	}()

	t.logger.Debugf("starting ffmpeg")
	t.logger.Debugf("%s", t.task.Cmdline)

	var stdout, stderr bytes.Buffer

	t.cmd.Stdout = &stdout
	t.cmd.Stderr = &stderr

	err := t.cmd.Start()
	if err != nil {
		fmtErr := fmt.Errorf("ffmpeg: %s\nstdout: %s\nstderr: %s\n", err, stdout.String(), stderr.String())
		errCh <- fmtErr
		return
	}

	t.logger.Info("transcoding")

	err = t.cmd.Wait()
	if err != nil {
		fmtErr := err
		if ErrExitStatusInterrupt.Error() == err.Error() {
			t.logger.Warning("ffmpeg execution has been canceled")
		} else {
			fmtErr = fmt.Errorf("ffmpeg: %s\nstdout: %s\nstderr: %s\n", err, stdout.String(), stderr.String())
		}

		errCh <- fmtErr
		return
	}
}
