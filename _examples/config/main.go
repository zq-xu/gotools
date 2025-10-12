package main

import (
	"github.com/zq-xu/gotools/config"

	"github.com/zq-xu/gotools"
)

var SampleConfig Config

type Config struct {
	Port  int
	Host  string
	Debug bool
}

func init() {
	config.Register("sample", &SampleConfig, func() error { return nil })
}

func main() {
	err := config.Setup("config.yaml")
	if err != nil {
		panic(err)
	}

	gotools.Logger.Info("Port:", SampleConfig.Port)
	gotools.Logger.Info("Host:", SampleConfig.Host)
	gotools.Logger.Info("Debug:", SampleConfig.Debug)
}
