package routerx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/routerx/auth"
)

var (
	AuthMiddleware        = auth.Middleware
	InitAuthMiddleware    = auth.InitAuthMiddleware
	GetUserInfoFromToken  = auth.GetUserInfoFromToken
	GetAccountInfoHandler = auth.GetAccountInfoHandler

	WithKey        = auth.WithKey
	WithTimeout    = auth.WithTimeout
	WithMaxRefresh = auth.WithMaxRefresh
)

func Login(ctx *gin.Context) {
	mw, err := auth.Middleware()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorx.ResponseBody{Error: err.Error()})
	}
	mw.LoginHandler(ctx)
}

func Logout(ctx *gin.Context) {
	mw, err := auth.Middleware()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorx.ResponseBody{Error: err.Error()})
	}
	mw.LogoutHandler(ctx)
}
