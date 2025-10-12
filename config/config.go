package config

import (
	"fmt"
	"sync"

	"github.com/rotisserie/eris"
	"github.com/spf13/viper"

	"github.com/zq-xu/gotools/logs"
)

var (
	registryMu sync.Mutex
	registry   = make([]setupItem, 0)
)

type setupFunc func() error

type setupItem struct {
	name string
	cfg  any
	fn   setupFunc
}

// Register
func Register(name string, cfg any, fn setupFunc) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = append(registry, setupItem{name, cfg, fn})
}

// Setup
func Setup(filename string) error {
	v := viper.New()
	v.SetConfigFile(filename)

	if err := v.ReadInConfig(); err != nil {
		return eris.Wrap(err, "read config file failed.")
	}

	registryMu.Lock()
	defer registryMu.Unlock()

	for _, item := range registry {
		err := loadConfig(v, item.name, item.cfg)
		if err != nil {
			return err
		}

		err = item.fn()
		if err != nil {
			return eris.Wrapf(err, "failed to setup %s", item.name)
		}
		logs.Logger.Infof("Succeed to init %s", item.name)
	}
	return nil
}

func loadConfig(v *viper.Viper, name string, value any) error {
	sub := v.Sub(name)
	if sub == nil {
		fmt.Println("empty config found for ", name)
		return nil
	}

	err := sub.Unmarshal(value)
	if err != nil {
		return eris.Wrapf(err, "unmarshal config for %s failed.", name)
	}

	return nil
}
