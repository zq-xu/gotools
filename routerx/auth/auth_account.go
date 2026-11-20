package auth

import (
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthAccount interface {
	GetID() string
	GetUsername() string
	GetRoles() string
	GetName() string
	GetStatus() string
}

type AccountInfoResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username" description:"the username for login"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	Status   string   `json:"status"`
}

// GetUserInfoHandler
func GetAccountInfoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, GetUserInfoFromToken(ctx))
}

// GetUserInfoFromToken 从 JWT claims 取当前用户信息。
// 只能用在已经过 JWT 中间件的路由上：未认证时 claims 为空，会返回全空的结构体
// （而不是 panic），调用方用 ID 是否为空判断即可。
func GetUserInfoFromToken(ctx *gin.Context) *AccountInfoResponse {
	claims := jwt.ExtractClaims(ctx)

	roles := claimString(claims, AuthUserRolesToken)

	return &AccountInfoResponse{
		ID:       claimString(claims, AuthAccountIDToken),
		Name:     claimString(claims, AuthAccountNameToken),
		Username: claimString(claims, AuthAccountUsernameToken),
		Roles:    strings.Split(roles, ","),
		Status:   claimString(claims, AuthAccountStatusToken),
	}
}

func claimString(claims jwt.MapClaims, key string) string {
	if v, ok := claims[key].(string); ok {
		return v
	}

	return ""
}
