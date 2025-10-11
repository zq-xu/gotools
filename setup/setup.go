package setup

import (
	"log"

	"github.com/rotisserie/eris"

	"zq-xu/gotools/config"
)

var setupSet = map[string]setupFunc{}

type setupFunc func(*config.Config) error

func RegisterSetup(name string, fn setupFunc) {
	setupSet[name] = fn
}

// Setup
func Setup() error {
	err := config.InitConfig()
	if err != nil {
		return eris.Wrap(err, "failed to init config")
	}

	return setup()
}

func setup() error {
	for name, fn := range setupSet {
		err := fn(config.Cfg)
		if err != nil {
			return eris.Wrapf(err, "failed to setup %s", name)
		}

		log.Println("Succeed to setup", name)
	}
	return nil
}
