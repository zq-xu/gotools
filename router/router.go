package router

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"

	"zq-xu/gotools/config"
	"zq-xu/gotools/utils"
)

const (
	// MaxMultipartMemory   // 100 * 2^20 = 100MB
	MaxMultipartMemory = 100 << 20
)

var groups = make([]*APIGroup, 0)

func init() {
	config.Register("router", &RouterConfig, func() error { return nil })
}

// RegisterGroup adds the route group into the route map
func RegisterGroup(grps ...*APIGroup) {
	for _, grp := range grps {
		if utils.IsInterfaceValueNil(grp) {
			return
		}

		groups = append(groups, grp)
	}
}

// StartRouter
func StartRouter(ctx context.Context, r *gin.Engine) error {
	srv := &http.Server{
		Addr:    net.JoinHostPort(RouterConfig.Host, fmt.Sprintf("%d", RouterConfig.Port)),
		Handler: r,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- startServe(srv)
	}()

	select {
	case <-ctx.Done():
		return shutdownServer(srv)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	}
}

func startServe(srv *http.Server) error {
	log.Println("Starting server at", srv.Addr)

	if RouterConfig.DisableTLS {
		return srv.ListenAndServe()
	}

	return srv.ListenAndServeTLS(RouterConfig.CertPath, RouterConfig.KeyPath)
}

func shutdownServer(srv *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("Server forced to shutdown.")

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		return eris.Wrap(err, "shutdown failed")
	}

	log.Println("Server exited gracefully")
	return nil
}

func NewRouter() *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.MaxMultipartMemory = MaxMultipartMemory

	r.Use(gin.Recovery())
	r.Use(loggerFilter([]string{HealthPath}))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(corsMiddleware())

	r.NoRoute(func(c *gin.Context) { c.JSON(404, gin.H{"api": "not found"}) })
	registerHealth(r)

	for _, grp := range groups {
		grp.AddToEngine(r)
	}

	return r
}
