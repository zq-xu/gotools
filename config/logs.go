package config

import (
	"github.com/zq-xu/gotools/logs"
)

func init() {
	Register("logs", &logs.LogConfig, logs.InitLogger)
}
