package fastapix

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/storex"
)

func GetHandler[T any, R any](ctx *gin.Context, transFn func(obj *T) (*R, errorx.ErrorInfo)) {
	id := routerx.GetID(ctx)

	resp, ei := get(ctx, id, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func get[T any, R any](ctx context.Context, id string,
	transFn func(obj *T) (*R, errorx.ErrorInfo)) (*R, errorx.ErrorInfo) {
	obj := new(T)
	ei := storex.NewErrorInfoForGetError(storex.DB(ctx).Get(obj, id))
	if ei != nil {
		return nil, ei
	}

	resp, ei := transFn(obj)
	if ei != nil {
		return nil, ei
	}

	return resp, nil
}
