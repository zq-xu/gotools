package auth

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthAccount interface {
	GetID() string
	GetUsername() string
	GetName() string
	GetStatus() string
}

type AccountInfoResponse struct {
	ID       string `json:"id"`
	Username string `json:"username" description:"the username for login"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

// GetUserInfoHandler
func GetAccountInfoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, getUserInfoFromToken(ctx))
}

func getUserInfoFromToken(ctx *gin.Context) *AccountInfoResponse {
	claims := jwt.ExtractClaims(ctx)

	return &AccountInfoResponse{
		ID:       claims[AuthAccountIDToken].(string),
		Name:     claims[AuthAccountNameToken].(string),
		Username: claims[AuthAccountUsernameToken].(string),
		Status:   claims[AuthAccountStatusToken].(string),
	}
}
