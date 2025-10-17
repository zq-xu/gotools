package fastapix

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/utilsx"
)

func GetByRawGormHandler[T any, R any](ctx *gin.Context,
	queryFn func(*gorm.DB, string) *gorm.DB,
	transFn func(obj *T) (*R, errorx.ErrorInfo)) {
	id := routerx.GetID(ctx)

	resp, ei := getByRawGorm(ctx, id, queryFn, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func getByRawGorm[T any, R any](ctx context.Context, id string, queryFn func(*gorm.DB, string) *gorm.DB, transFn func(obj *T) (*R, errorx.ErrorInfo)) (*R, errorx.ErrorInfo) {
	obj := new(T)

	err := queryFn(storex.GormDB(ctx), id).First(obj).Error
	ei := storex.NewErrorInfoForGetError(err)
	if ei != nil {
		return nil, ei
	}

	return transFn(obj)
}

func DefaultTransObjToResp[T any, R any](obj *T) (*R, errorx.ErrorInfo) {
	resp := new(R)

	err := utilsx.Copy(resp, obj)
	if err != nil {
		return nil, errorx.NewError(http.StatusBadRequest, "failed to copy", err)
	}
	return resp, nil
}
