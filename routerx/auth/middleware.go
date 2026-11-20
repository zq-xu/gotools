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

	"github.com/zq-xu/gotools/errorx"
)

const (
	AuthAccountIDToken       = "auth_account_id"
	AuthAccountNameToken     = "auth_account_name"
	AuthAccountUsernameToken = "auth_account_username"
	AuthUserRolesToken       = "auth_account_roles"
	AuthAccountStatusToken   = "auth_account_status"
)

var (
	middleware    *jwt.GinJWTMiddleware
	accountLoader func(ctx context.Context, username, password string) (AuthAccount, errorx.ErrorInfo)
)

// Options JWT 中间件可配置项，零值时使用默认值。
type Options struct {
	// Key 签名密钥，为空时使用不安全的默认值（仅限开发环境）。
	Key []byte
	// Timeout token 有效期，为 0 时默认 1 小时。
	Timeout time.Duration
	// MaxRefresh 可刷新窗口，为 0 时默认 100 小时。
	MaxRefresh time.Duration
}

type Option func(*Options)

// WithKey 指定 JWT 签名密钥。
func WithKey(key string) Option {
	return func(o *Options) {
		o.Key = []byte(key)
	}
}

// WithTimeout 指定 token 有效期。
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

// WithMaxRefresh 指定 token 可刷新窗口（需不小于 Timeout）。
func WithMaxRefresh(d time.Duration) Option {
	return func(o *Options) {
		o.MaxRefresh = d
	}
}

const (
	defaultKey        = "secret key"
	defaultTimeout    = time.Hour
	defaultMaxRefresh = time.Hour * 100
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

// InitAuthMiddleware 初始化 JWT 中间件，可用 Option 覆盖 Key / Timeout / MaxRefresh。
func InitAuthMiddleware(loader func(ctx context.Context, username, password string) (AuthAccount, errorx.ErrorInfo), opts ...Option) {
	var err error
	accountLoader = loader

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	if len(o.Key) == 0 {
		o.Key = []byte(defaultKey)
		log.Printf("routerx/auth: using default JWT key, use auth.WithKey to set a real one in production")
	}
	if o.Timeout <= 0 {
		o.Timeout = defaultTimeout
	}
	if o.MaxRefresh <= 0 {
		o.MaxRefresh = defaultMaxRefresh
	}

	middleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Key:             o.Key,
		Timeout:         o.Timeout,
		MaxRefresh:      o.MaxRefresh,
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

func validateAccount(ctx *gin.Context, u *LoginAccount) errorx.ErrorInfo {
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
			AuthUserRolesToken:       v.authAccount.GetRoles(),
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
	ctx.Set(AuthUserRolesToken, claims[AuthUserRolesToken].(string))
	ctx.Set(AuthAccountStatusToken, claims[AuthAccountStatusToken].(string))
	return nil
}

func unauthorized(ctx *gin.Context, code int, message string) {
	if strings.Contains(strings.ToLower(message), "expired") {
		ctx.JSON(http.StatusUnauthorized, &UnauthorizedResp{
			Code:    http.StatusUnauthorized,
			Message: "token expired",
		})
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
