package gormkit

import (
	"github.com/zq-xu/gotools/configx"
	"gorm.io/gorm/logger"
)

var GormConfig Config

type Config struct {
	Address      string `yaml:"address"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database"`

	LogLevel logger.LogLevel `yaml:"logLevel"`
}

func init() {
	configx.Register("database", &GormConfig, InitGorm)
}
