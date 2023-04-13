package auth

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/helper/restapi/response"
	"zq-xu/helper/store"
)

const (
	AuthUserIDToken = "auth_user_id"
)

var (
	Middleware *jwt.GinJWTMiddleware

	gjm = &jwt.GinJWTMiddleware{
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     "user",
		PayloadFunc:     generatePayLoad,
		IdentityHandler: identityHandler,
		Authenticator:   authenticate,
		Authorizator:    authorize,
		LoginResponse:   loginResponse,
		RefreshResponse: loginResponse,
		LogoutResponse:  logoutResponse,
		Unauthorized:    unauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
)

type UserReq struct {
	ID       string
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

type UnauthorizedResp struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

func init() {
	Middleware, _ = jwt.New(gjm)
}

func generatePayLoad(data interface{}) jwt.MapClaims {
	if v, ok := data.(*UserReq); ok {
		return jwt.MapClaims{
			AuthUserIDToken: v.ID,
		}
	}

	return jwt.MapClaims{}
}

func identityHandler(ctx *gin.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	u := &UserReq{
		ID: claims[AuthUserIDToken].(string),
	}

	ctx.Set(AuthUserIDToken, u.ID)
	return u
}

func authenticate(ctx *gin.Context) (interface{}, error) {
	user := &UserReq{}
	if err := ctx.ShouldBind(user); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	err := validateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func validateUser(ctx *gin.Context, u *UserReq) error {
	db := store.DB(ctx)

	mu := &User{}
	err := db.
		Where("name = ? AND password = ?", u.Username, u.Password).
		First(mu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return jwt.ErrFailedAuthentication
		}
		return response.NewStorageError(response.StorageErrorCode, err)
	}

	u.ID = strconv.FormatInt(mu.ID, 10)
	return nil
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

func authorize(data interface{}, ctx *gin.Context) bool {
	//v, ok := data.(*User)
	//if !ok {
	//	return false
	//}
	//
	//err := validateUser(v)
	//if err != nil {
	//	return false
	//}

	return true
}

func unauthorized(ctx *gin.Context, code int, message string) {
	ctx.JSON(code,
		&UnauthorizedResp{
			Code:    code,
			Message: message,
		},
	)
}

func GetAuthUserId(ctx *gin.Context) string {
	return ctx.GetString(AuthUserIDToken)
}
