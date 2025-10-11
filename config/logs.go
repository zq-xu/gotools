package config

import (
	"zq-xu/gotools/logs"
)

func init() {
	Register("logs", &logs.LogConfig, logs.InitLogger)
}
