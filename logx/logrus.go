package logx

import (
	"os"

	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
)

const (
	DefaultLogrusLogLevel = logrus.InfoLevel
)

var (
	Logger = logrus.New()
)

var LogConfig Config

type Config struct {
	LogLevel string
}

func init() {
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.TraceLevel)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&MyFormatter{})

}

// InitLogger
func InitLogger() error {
	return SetLoggerLevel(LogConfig.LogLevel)
}

// SetLoggerLevel
func SetLoggerLevel(l string) error {
	level, err := getLogrusLevel(l)
	if err != nil {
		return err
	}
	Logger.SetLevel(level)
	return nil
}

func getLogrusLevel(str string) (logrus.Level, error) {
	if str == "" {
		return DefaultLogrusLogLevel, nil
	}

	l, err := logrus.ParseLevel(str)
	if err != nil {
		return DefaultLogrusLogLevel, eris.Wrapf(err, "failed to parse logrus level from %s", str)

	}

	return l, nil
}
