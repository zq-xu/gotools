package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"

	"zq-xu/gotools/apperror"
)

const (
	AuthAccountIDToken       = "auth_account_id"
	AuthAccountNameToken     = "auth_account_name"
	AuthAccountUsernameToken = "auth_account_username"
	AuthAccountStatusToken   = "auth_account_status"
)

var (
	middleware    *jwt.GinJWTMiddleware
	accountLoader func(ctx context.Context, username, password string) (AuthAccount, apperror.ErrorInfo)
)

type LoginAccount struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`

	authAccount AuthAccount
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

type UnauthorizedResp struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

func Middleware() (*jwt.GinJWTMiddleware, error) {
	if middleware == nil {
		return nil, eris.New("empty middleware")
	}
	return middleware, nil
}

// InitMiddleware
func InitMiddleware(loader func(ctx context.Context, username, password string) (AuthAccount, apperror.ErrorInfo)) {
	var err error
	accountLoader = loader

	middleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour * 100,
		IdentityKey:     "account",
		PayloadFunc:     generatePayLoad,
		IdentityHandler: identityHandler,
		Authenticator:   authenticate,
		LoginResponse:   loginResponse,
		RefreshResponse: loginResponse,
		LogoutResponse:  logoutResponse,
		Unauthorized:    unauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatalf("failed to initial auth middleware")
	}
}

func authenticate(ctx *gin.Context) (interface{}, error) {
	account := &LoginAccount{}

	err := ctx.ShouldBind(account)
	if err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	err = validateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func validateAccount(ctx *gin.Context, u *LoginAccount) apperror.ErrorInfo {
	ui, ei := accountLoader(ctx, u.Username, u.Password)
	if ei != nil {
		return ei
	}

	u.authAccount = ui
	return nil
}

func generatePayLoad(data interface{}) jwt.MapClaims {
	if v, ok := data.(*LoginAccount); ok {
		return jwt.MapClaims{
			AuthAccountIDToken:       v.authAccount.GetID(),
			AuthAccountNameToken:     v.authAccount.GetName(),
			AuthAccountUsernameToken: v.authAccount.GetUsername(),
			AuthAccountStatusToken:   v.authAccount.GetStatus(),
		}
	}

	return jwt.MapClaims{}
}

func identityHandler(ctx *gin.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	ctx.Set(AuthAccountIDToken, claims[AuthAccountIDToken].(string))
	ctx.Set(AuthAccountNameToken, claims[AuthAccountNameToken].(string))
	ctx.Set(AuthAccountUsernameToken, claims[AuthAccountUsernameToken].(string))
	ctx.Set(AuthAccountStatusToken, claims[AuthAccountStatusToken].(string))
	return nil
}

func unauthorized(ctx *gin.Context, code int, message string) {
	if strings.Contains(strings.ToLower(message), "expired") {
		ctx.JSON(http.StatusUnauthorized, "token expired")
		return
	}

	ctx.JSON(code,
		&UnauthorizedResp{
			Code:    code,
			Message: message,
		},
	)
}

func loginResponse(ctx *gin.Context, code int, token string, expire time.Time) {
	ctx.JSON(http.StatusCreated,
		&LoginResp{
			Token:  token,
			Expire: expire.Format(time.RFC3339),
		})
}

func logoutResponse(ctx *gin.Context, code int) {
	ctx.JSON(http.StatusCreated, struct{}{})
}
