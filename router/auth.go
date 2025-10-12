package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/router/auth"
)

var (
	AuthMiddleware        = auth.Middleware
	InitAuthMiddleware    = auth.InitAuthMiddleware
	GetUserInfoFromToken  = auth.GetUserInfoFromToken
	GetAccountInfoHandler = auth.GetAccountInfoHandler
)

func Login(ctx *gin.Context) {
	mw, err := auth.Middleware()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	mw.LoginHandler(ctx)
}

func Logout(ctx *gin.Context) {
	mw, err := auth.Middleware()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	mw.LogoutHandler(ctx)
}
