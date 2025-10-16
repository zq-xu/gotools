package fastapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/store/database"
	"github.com/zq-xu/gotools/types"
)

type listOpts func(*types.ListParams)

func WithFilter(key string, value string) func(*gotools.ListParams) {
	return func(listParam *types.ListParams) {
		if listParam.Queries == nil {
			listParam.Queries = make(types.Queries)
		}

		listParam.Queries[key] = value
	}
}

func ListHandler[T any](ctx *gin.Context,
	queryFn func(db database.Database, listParam *gotools.ListParams) ([]T, int, gotools.ErrorInfo),
	transFn func(count int, listParam *gotools.ListParams, listObj []T) (*gotools.PageResponse, gotools.ErrorInfo),
	opts ...listOpts) {
	listParam, ei := types.GetListParams(ctx)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	for _, opt := range opts {
		opt(listParam)
	}

	pageInfo, ei := list(ctx, listParam, queryFn, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, pageInfo)
}

func list[T any](ctx context.Context, listParam *types.ListParams,
	queryFn func(db database.Database, listParam *types.ListParams) ([]T, int, apperror.ErrorInfo),
	transFn func(count int, listParam *types.ListParams, listObj []T) (*types.PageResponse, apperror.ErrorInfo)) (*types.PageResponse, apperror.ErrorInfo) {
	db := store.DB(ctx)

	listObj, count, ei := queryFn(db, listParam)
	if ei != nil {
		return nil, ei
	}

	return transFn(int(count), listParam, listObj)
}

func DefaultListObjWithCount[T any](db database.Database, listParam *types.ListParams) ([]T, int, apperror.ErrorInfo) {
	listObj := make([]T, 0)
	count, err := db.ListWithCount(listParam, new(T), &listObj)
	if err != nil {
		return nil, 0, apperror.NewError(http.StatusBadRequest, "failed to list", err)
	}
	return listObj, int(count), nil
}
