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

func GetUserInfoFromToken(ctx *gin.Context) *AccountInfoResponse {
	claims := jwt.ExtractClaims(ctx)

	return &AccountInfoResponse{
		ID:       claims[AuthAccountIDToken].(string),
		Name:     claims[AuthAccountNameToken].(string),
		Username: claims[AuthAccountUsernameToken].(string),
		Roles:    strings.Split(claims[AuthUserRolesToken].(string), ","),
		Status:   claims[AuthAccountStatusToken].(string),
	}
}
