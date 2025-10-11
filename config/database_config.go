package config

import (
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Address      string `yaml:"address"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database"`

	LogLevel logger.LogLevel `yaml:"logLevel"`
}
