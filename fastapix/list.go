package fastapix

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/storex/database"
	"github.com/zq-xu/gotools/typesx"
)

type listOpts func(*typesx.ListParams)

func WithFilter(key string, value string) func(*typesx.ListParams) {
	return func(listParam *typesx.ListParams) {
		if listParam.Queries == nil {
			listParam.Queries = make(typesx.Queries)
		}

		listParam.Queries[key] = value
	}
}

func ListHandler[T any](ctx *gin.Context,
	queryFn func(db database.Database, listParam *typesx.ListParams) ([]T, int, errorx.ErrorInfo),
	transFn func(count int, listParam *typesx.ListParams, listObj []T) (*typesx.PageResponse, errorx.ErrorInfo),
	opts ...listOpts) {
	listParam, ei := typesx.GetListParams(ctx)
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

func list[T any](ctx context.Context, listParam *typesx.ListParams,
	queryFn func(db database.Database, listParam *typesx.ListParams) ([]T, int, errorx.ErrorInfo),
	transFn func(count int, listParam *typesx.ListParams, listObj []T) (*typesx.PageResponse, errorx.ErrorInfo)) (*typesx.PageResponse, errorx.ErrorInfo) {
	db := storex.DB(ctx)

	listObj, count, ei := queryFn(db, listParam)
	if ei != nil {
		return nil, ei
	}

	return transFn(int(count), listParam, listObj)
}

func DefaultListObjWithCount[T any](db database.Database, listParam *typesx.ListParams) ([]T, int, errorx.ErrorInfo) {
	listObj := make([]T, 0)
	count, err := db.ListAssociationsWithCount(listParam, new(T), &listObj)
	if err != nil {
		return nil, 0, errorx.NewError(http.StatusBadRequest, "failed to list", err)
	}
	return listObj, int(count), nil
}
