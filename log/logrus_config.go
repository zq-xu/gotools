package log

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"zq-xu/helper/config"
)

const (
	logrusConfigName = "LogrusConfig"

	logrusLogLevelEnv     = "LogLevel"
	defaultLogrusLogLevel = logrus.InfoLevel
)

type LogrusConfig struct {
	LogLevel string

	Level logrus.Level
}

var (
	LogrusCfg = &LogrusConfig{}
)

func init() {
	config.RegisterCfg(logrusConfigName, LogrusCfg)
}

// AddFlags adds flags for logrus
func (lc *LogrusConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&lc.LogLevel, "log-level", os.Getenv(logrusLogLevelEnv), "the log level")
}

func (lc *LogrusConfig) Revise() {
	lc.reviseLevel()
}

func (lc *LogrusConfig) reviseLevel() {
	if lc.LogLevel == "" {
		lc.Level = defaultLogrusLogLevel
		return
	}

	l, err := logrus.ParseLevel(lc.LogLevel)
	if err != nil {
		Logger.Warningf("Failed to parse the log level! %v", err)
		lc.Level = defaultLogrusLogLevel
		return
	}

	lc.Level = l
}
