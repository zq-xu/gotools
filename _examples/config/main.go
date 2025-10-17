package main

import (
	"github.com/zq-xu/gotools/configx"
	"github.com/zq-xu/gotools/logx"
)

var SampleConfig Config

type Config struct {
	Port  int
	Host  string
	Debug bool
}

func init() {
	configx.Register("sample", &SampleConfig, func() error { return nil })
}

func main() {
	err := configx.Setup("configx.yaml")
	if err != nil {
		panic(err)
	}

	logx.Logger.Info("Port:", SampleConfig.Port)
	logx.Logger.Info("Host:", SampleConfig.Host)
	logx.Logger.Info("Debug:", SampleConfig.Debug)
}
