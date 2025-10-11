package config

import (
	"os"

	"github.com/rotisserie/eris"
	"sigs.k8s.io/yaml"

	"zq-xu/gotools/utils"
)

const (
	defaultConfigFilePath = "./config.yaml"
	configFilePathEnvKey  = "CONFIG_FILE_PATH"
)

var (
	Cfg = &Config{}
)

type Config struct {
	LogLevel string `yaml:"logLevel"`
	AesKey   string `yaml:"aesKey"`

	DatabaseConfig DatabaseConfig `yaml:"databaseConfig"`
	RouteConfig    RouteConfig    `yaml:"routeConfig"`
}

// InitConfig
func InitConfig() error {
	return initConfig(Cfg)
}

func InitCustomisedConfig(i interface{}) error {
	return initConfig(i)
}

func initConfig(i interface{}) error {
	filePath := os.Getenv(configFilePathEnvKey)
	if filePath == "" {
		filePath = defaultConfigFilePath
	}

	bs, err := utils.ReadFiles(filePath)
	if err != nil {
		return eris.Wrapf(err, "failed to read config file %s", filePath)
	}

	return yaml.Unmarshal(bs, i)
}
