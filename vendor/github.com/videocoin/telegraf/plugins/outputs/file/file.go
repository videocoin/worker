package file

import (
	"fmt"
	"io"
	"os"

	"github.com/videocoin/telegraf"
	logger "github.com/videocoin/telegraf/log"
	"github.com/videocoin/telegraf/plugins/outputs"
	"github.com/videocoin/telegraf/plugins/serializers"
	"github.com/sirupsen/logrus"
)

type File struct {
	Files []string

	writers []io.Writer
	closers []io.Closer

	serializer serializers.Serializer
}

var log *logrus.Entry

func init() {
	log = logger.Get()
}

func (f *File) SetSerializer(serializer serializers.Serializer) {
	f.serializer = serializer
}

func (f *File) Connect() error {
	if len(f.Files) == 0 {
		f.Files = []string{"stdout"}
	}

	for _, file := range f.Files {
		if file == "stdout" {
			f.writers = append(f.writers, os.Stdout)
		} else {
			of, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend|0644)
			if err != nil {
				return err
			}

			f.writers = append(f.writers, of)
			f.closers = append(f.closers, of)
		}
	}
	return nil
}

func (f *File) Close() error {
	var err error
	for _, c := range f.closers {
		errClose := c.Close()
		if errClose != nil {
			err = errClose
		}
	}
	return err
}

func (f *File) Write(metrics []telegraf.Metric) error {
	var writeErr error = nil
	for _, metric := range metrics {
		b, err := f.serializer.Serialize(metric)
		if err != nil {
			return fmt.Errorf("failed to serialize message: %s", err)
		}

		for _, writer := range f.writers {
			_, err = writer.Write(b)
			if err != nil && writer != os.Stdout {
				writeErr = fmt.Errorf("failed to write message: %s, %s", b, err)
			}
		}
	}
	return writeErr
}

func init() {
	outputs.Add("file", func() telegraf.Output {
		return &File{}
	})
}
