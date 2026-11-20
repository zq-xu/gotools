package getx

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/utilsx"
)

type validateFn func(ctx *gin.Context) errorx.ErrorInfo

// ValidateAuthInRoles 要求 token 里的 roles 至少命中传入的其中一个，否则返回 403。
// 它只能挂在已经过 JWT 中间件的路由上：claims 是从 ctx 里取的，没有中间件就是空的。
func ValidateAuthInRoles(roles ...string) validateFn {
	return func(ctx *gin.Context) errorx.ErrorInfo {
		userInfo := routerx.GetUserInfoFromToken(ctx)
		if userInfo == nil || userInfo.ID == "" {
			return errorx.NewError(http.StatusUnauthorized, "missing user info in token", nil)
		}

		for _, r := range userInfo.Roles {
			if slices.Contains(roles, r) {
				return nil
			}
		}

		return errorx.NewError(http.StatusForbidden, "permission denied", nil)
	}
}

func DefaultTransObjToResp[T any, R any](obj *T) (*R, errorx.ErrorInfo) {
	resp := new(R)

	err := utilsx.Copy(resp, obj)
	if err != nil {
		return nil, errorx.NewError(http.StatusBadRequest, "failed to copy", err)
	}
	return resp, nil
}
