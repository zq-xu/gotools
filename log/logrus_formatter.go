package log

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const timeFormatter = "2006-01-02 15:04:05"

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := entry.Buffer

	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(timeFormatter)
	if entry.HasCaller() {
		//b.WriteString(fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
		//	timestamp, entry.Level,
		//	filepath.Base(entry.Caller.File), entry.Caller.Line, entry.Caller.Function,
		//	entry.Message))
		//return b.Bytes(), nil

		b.WriteString(fmt.Sprintf("[%s] [%s] [%s:%d] %s\n",
			timestamp, entry.Level,
			filepath.Base(entry.Caller.File), entry.Caller.Line,
			entry.Message))
		return b.Bytes(), nil
	}

	b.WriteString(fmt.Sprintf("[%s] [%s] %s\n",
		timestamp, entry.Level,
		entry.Message))
	return b.Bytes(), nil
}
