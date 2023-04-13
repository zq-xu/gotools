package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	HealthPath = "/health"
)

var (
	healthFunc = Health
)

func SetHealthHandler(healthHandler func(ctx *gin.Context)) {
	healthFunc = healthHandler
}

func registerHealth(r *gin.Engine) {
	r.GET(HealthPath, healthFunc)
}

func Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &struct{}{})
}
