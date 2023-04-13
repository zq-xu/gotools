package router

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"

	"zq-xu/helper/log"
)

var pprofOnce = &sync.Once{}

func StartPprof() {
	pprofOnce.Do(func() {
		addr := fmt.Sprintf("%s:%s", RouteCfg.IP, RouteCfg.PprofPort)

		log.Logger.Infof("start to listen the pprof on %v", addr)

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Logger.Infof("start pprof failed on %s", addr)
			os.Exit(1)
		}
	})
}
