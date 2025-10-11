package gormkit

import (
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
