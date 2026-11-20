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

func ValidateAuthInRoles(roles ...string) validateFn {
	return func(ctx *gin.Context) errorx.ErrorInfo {
		userInfo := routerx.GetUserInfoFromToken(ctx)
		for _, v := range roles {
			if slices.Contains(userInfo.Roles, v) {
				return nil
			}
		}
		return nil
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
