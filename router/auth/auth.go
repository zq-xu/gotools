package auth

import (
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	Middleware.LoginHandler(ctx)
}

func Logout(ctx *gin.Context) {
	Middleware.LogoutHandler(ctx)
}

func RefreshToken(ctx *gin.Context) {
	Middleware.RefreshHandler(ctx)
}
