package router

import (
	"net"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"zq-xu/helper/utils"
)

var groups = make([]*APIGroup, 0)

type engineOpt func(r *gin.Engine)

func RegisterGroup(grps ...*APIGroup) {
	for _, grp := range grps {
		if utils.IsInterfaceValueNil(grp) {
			return
		}

		groups = append(groups, grp)
	}
}

func StartRouter(r *gin.Engine) error {
	addr := net.JoinHostPort(RouteCfg.IP, RouteCfg.Port)

	for _, grp := range groups {
		grp.AddToEngine(r)
	}

	if RouteCfg.DisableTLS {
		return r.Run(addr)
	}

	return r.RunTLS(addr, RouteCfg.CertPath, RouteCfg.KeyPath)
}

func NewRouter(fns ...engineOpt) *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	for _, fn := range fns {
		fn(r)
	}

	return r
}

func DefaultRouter() *gin.Engine {
	return NewRouter(func(r *gin.Engine) {
		r.Use(LoggerFilter([]string{HealthPath}, GetMethodFilter))
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		r.Use(corsMiddleware())
		registerHealth(r)
	})
}
