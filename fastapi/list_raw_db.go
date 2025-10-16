package fastapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/store"
	"gorm.io/gorm"
)

func ListByRawGormHandler[T any](ctx *gin.Context,
	queryFn func(db *gorm.DB, listParam *gotools.ListParams) *gorm.DB,
	transFn func(count int, listParam *gotools.ListParams, listObj []T) (*gotools.PageResponse, gotools.ErrorInfo)) {
	listParam, ei := gotools.GetListParams(ctx)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	resp, ei := listByRawGorm(ctx, listParam, queryFn, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func listByRawGorm[T any](ctx context.Context, listParam *gotools.ListParams,
	queryFn func(db *gorm.DB, listParam *gotools.ListParams) *gorm.DB,
	transFn func(count int, listParam *gotools.ListParams, listObj []T) (*gotools.PageResponse, gotools.ErrorInfo)) (*gotools.PageResponse, gotools.ErrorInfo) {
	count, err := store.DB(ctx).GetCount(new(T), listParam)
	if err != nil {
		return nil, gotools.NewError(http.StatusBadRequest, "get count failed", err)
	}

	listObj := make([]T, 0)
	db := queryFn(store.GormDB(ctx), listParam)
	db = store.GenerateDBForQuery(db, listParam, &listObj)
	db = store.OptFuzzySearchDB(db, nil, listParam.FuzzySearchValue)
	db = store.OptPageDB(db, listParam)

	if err := db.Find(&listObj).Error; err != nil {
		return nil, gotools.NewError(http.StatusBadRequest, "load models failed", err)
	}

	return transFn(int(count), listParam, listObj)
}

func DefaultTransListObjToResp[T any, R any](count int, listParam *gotools.ListParams, listObj []T) (*gotools.PageResponse, gotools.ErrorInfo) {
	items := make([]interface{}, 0)

	for _, v := range listObj {
		r := new(R)

		err := gotools.Copy(&r, &v)
		if err != nil {
			return nil, gotools.NewError(http.StatusBadRequest, "failed to copy", err)
		}

		items = append(items, r)
	}

	return gotools.NewPageResponse(count, listParam.PageInfo, items), nil
}
