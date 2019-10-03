package transcoder

import "errors"

var (
	ErrExitStatusInterrupt    = errors.New("exit status 255")
	ErrProcessAlreadyFinished = errors.New("os: process already finished")
)
