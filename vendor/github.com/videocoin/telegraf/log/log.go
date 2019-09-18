package log

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var log *logrus.Entry

func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	log = logrus.NewEntry(logger).WithField("system", "stats")
}

func Get() *logrus.Entry {
	return log
}
