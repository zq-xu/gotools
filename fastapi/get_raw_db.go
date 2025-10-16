package fastapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/router"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/utils"
)

func GetByRawGormHandler[T any, R any](ctx *gin.Context,
	queryFn func(*gorm.DB, string) *gorm.DB,
	transFn func(obj *T) (*R, gotools.ErrorInfo)) {
	id := router.GetID(ctx)

	resp, ei := getByRawGorm(ctx, id, queryFn, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func getByRawGorm[T any, R any](ctx context.Context, id string, queryFn func(*gorm.DB, string) *gorm.DB, transFn func(obj *T) (*R, apperror.ErrorInfo)) (*R, apperror.ErrorInfo) {
	obj := new(T)

	err := queryFn(store.GormDB(ctx), id).First(obj).Error
	ei := store.NewErrorInfoForGetError(err)
	if ei != nil {
		return nil, ei
	}

	return transFn(obj)
}

func DefaultTransObjToResp[T any, R any](obj *T) (*R, apperror.ErrorInfo) {
	resp := new(R)

	err := utils.Copy(resp, obj)
	if err != nil {
		return nil, apperror.NewError(http.StatusBadRequest, "failed to copy", err)
	}
	return resp, nil
}
