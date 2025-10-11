package logs

import (
	"os"

	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"

	"zq-xu/gotools/config"
)

const (
	defaultLogrusLogLevel = logrus.InfoLevel
)

var Logger = logrus.New()

func init() {
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.TraceLevel)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&MyFormatter{})
}

// InitLogger
func InitLogger(cfg *config.Config) error {
	level, err := getLogrusLevel(cfg.LogLevel)
	if err != nil {
		return err
	}
	Logger.SetLevel(level)

	Logger.Infof("Succeed to init log with level %s", Logger.Level.String())
	return nil
}

func getLogrusLevel(str string) (logrus.Level, error) {
	if str == "" {
		return defaultLogrusLogLevel, nil
	}

	l, err := logrus.ParseLevel(str)
	if err != nil {
		return defaultLogrusLogLevel, eris.Wrapf(err, "failed to parse logrus level from %s", str)

	}

	return l, nil
}
