package transcoder

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/armon/circbuf"
	"github.com/sirupsen/logrus"
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

	var (
		stdoutPipe, stderrPipe io.ReadCloser
		err                    error
	)

	err = nil

	if t.logger.Logger.GetLevel() == logrus.DebugLevel {
		stdoutPipe, err = t.cmd.StdoutPipe()
		if err != nil {
			fmtErr := fmt.Errorf("ffmpeg: %s", err)
			errCh <- fmtErr
			return
		}

		stderrPipe, err = t.cmd.StderrPipe()
		if err != nil {
			fmtErr := fmt.Errorf("ffmpeg: %s", err)
			errCh <- fmtErr
			return
		}
	}

	err = t.cmd.Start()
	if err != nil {
		fmtErr := fmt.Errorf("ffmpeg: %s", err)
		errCh <- fmtErr
		return
	}

	if t.logger.Logger.GetLevel() == logrus.DebugLevel {
		stdouterr := bytes.NewBuffer(nil)

		go func() {
			_, err := io.Copy(stdouterr, stderrPipe)
			if err != nil {
				t.logger.Error(err)
			}
		}()
		go func() {
			_, err := io.Copy(stdouterr, stdoutPipe)
			if err != nil {
				t.logger.Error(err)
			}
		}()

		go func(stopCh chan bool) {
			ffmpegout, _ := circbuf.NewBuffer(1024 * 4)
			for {
				select {
				case <-stopCh:
					return
				default:
					buf := bytes.NewBuffer(nil)
					for {
						b, err := stdouterr.ReadByte()
						if err != nil {
							break
						}

						if b == '\x00' || b == 'x' {
							continue
						}

						if b == '\r' || b == '\n' {
							if ffmpegout != nil {
								_, err = ffmpegout.Write([]byte{'\n'})
								if err != nil {
									break
								}
							}

							line := buf.String()
							buf.Reset()

							if len(line) > 0 {
								t.logger.
									WithField("system", "ffmpeg").
									Debugf("ffmpeg: %s", line)
							}
						} else {
							if ffmpegout != nil {
								_, err = ffmpegout.Write([]byte{b})
								if err != nil {
									break
								}
							}
							buf.WriteByte(b)
						}
					}
				}
			}
		}(stopCh)
	}

	t.logger.Info("transcoding")

	err = t.cmd.Wait()
	if err != nil {
		fmtErr := err
		if ErrExitStatusInterrupt.Error() == err.Error() {
			t.logger.Warning("ffmpeg execution has been canceled")
		} else {
			fmtErr = fmt.Errorf("ffmpeg: %s", err)
		}

		errCh <- fmtErr
		return
	}
}
