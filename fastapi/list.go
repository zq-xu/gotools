package fastapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/types"
)

type listOpts func(*gotools.ListParams)

func WithFilter(key string, value string) func(*gotools.ListParams) {
	return func(listParam *gotools.ListParams) {
		if listParam.Queries == nil {
			listParam.Queries = make(types.Queries)
		}

		listParam.Queries[key] = value
	}
}

func ListHandler[T any](ctx *gin.Context,
	transFn func(count int, listParam *gotools.ListParams, listObj []T) (*gotools.PageResponse, gotools.ErrorInfo),
	opts ...listOpts) {
	listParam, ei := gotools.GetListParams(ctx)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	for _, opt := range opts {
		opt(listParam)
	}

	pageInfo, ei := list[T](ctx, listParam, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, pageInfo)
}

func list[T any](ctx context.Context, listParam *gotools.ListParams,
	transFn func(count int, listParam *gotools.ListParams, listObj []T) (*gotools.PageResponse, gotools.ErrorInfo)) (*gotools.PageResponse, gotools.ErrorInfo) {
	listObj := make([]T, 0)
	count, err := store.DB(ctx).ListWithCount(listParam, new(T), &listObj)
	if err != nil {
		return nil, gotools.NewError(http.StatusBadRequest, "load models failed", err)
	}

	return transFn(int(count), listParam, listObj)
}
