package fastapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/router"
	"github.com/zq-xu/gotools/store"
)

func GetHandler[T any, R any](ctx *gin.Context, transFn func(obj *T) (*R, gotools.ErrorInfo)) {
	id := router.GetID(ctx)

	resp, ei := get(ctx, id, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func get[T any, R any](ctx context.Context, id string,
	transFn func(obj *T) (*R, apperror.ErrorInfo)) (*R, apperror.ErrorInfo) {
	obj := new(T)
	ei := store.NewErrorInfoForGetError(store.DB(ctx).Get(obj, id))
	if ei != nil {
		return nil, ei
	}

	resp, ei := transFn(obj)
	if ei != nil {
		return nil, ei
	}

	return resp, nil
}
