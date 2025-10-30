package configx

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/rotisserie/eris"
	"github.com/spf13/viper"

	"github.com/zq-xu/gotools/logx"
)

var (
	registryMu sync.Mutex
	registry   = make([]*setupItem, 0)
)

type setupFunc func() error

type setupItem struct {
	name   string
	cfg    any
	fn     setupFunc
	byFile bool
}

func init() {
	Register("logs", &logx.LogConfig, logx.InitLogger)
}

// Register
func Register(name string, cfg any, fn setupFunc) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = append(registry, &setupItem{name, cfg, fn, false})
}

// Register
func RegisterByFile(name string, cfg any, fn setupFunc) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = append(registry, &setupItem{name, cfg, fn, true})
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
		err := loadConfig(v, item)
		if err != nil {
			return err
		}

		err = item.fn()
		if err != nil {
			return eris.Wrapf(err, "failed to setup %s", item.name)
		}
		logx.Logger.Infof("Succeed to init %s", item.name)
	}
	return nil
}

func loadConfig(v *viper.Viper, item *setupItem) error {
	if item.byFile {
		return loadConfigFromFilePath(v, item)
	}

	return loadConfigFromSub(v, item)
}

func loadConfigFromFilePath(v *viper.Viper, item *setupItem) error {
	filePath := v.GetString(item.name)

	vFile := viper.New()
	vFile.SetConfigFile(filePath)

	if err := vFile.ReadInConfig(); err != nil {
		return eris.Wrapf(err, "failed to read config file %s: %s", item.name, filePath)
	}

	if err := vFile.Unmarshal(item.cfg); err != nil {
		return eris.Wrapf(err, "failed to unmarshal config from file: %s", filePath)
	}

	return nil
}

func loadConfigFromSub(v *viper.Viper, item *setupItem) error {
	sub := v.Sub(item.name)
	if sub == nil {
		fmt.Println("empty config found for", item.name)
		return nil
	}

	if err := sub.Unmarshal(item.cfg); err != nil {
		return eris.Wrapf(err, "failed to unmarshal config for: %s", item.name)
	}

	return nil
}

// DefaultSetupFunc
func DefaultSetupFunc() error { return nil }

// DebugSetupFunc
func DebugSetupFunc(k string, v any) setupFunc {
	return func() error {
		b, _ := json.Marshal(v)
		logx.Logger.Infof("%s: %+v", k, string(b))
		return nil
	}
}
