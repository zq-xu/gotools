package fastapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/types"
	"github.com/zq-xu/gotools/utils"
)

func ListByRawGormHandler[T any](ctx *gin.Context,
	queryFn func(db *gorm.DB, listParam *gotools.ListParams) *gorm.DB,
	transFn func(count int, listParam *gotools.ListParams, listObj []T) (*types.PageResponse, apperror.ErrorInfo)) {
	listParam, ei := types.GetListParams(ctx)
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

func listByRawGorm[T any](ctx context.Context, listParam *types.ListParams,
	queryFn func(db *gorm.DB, listParam *types.ListParams) *gorm.DB,
	transFn func(count int, listParam *types.ListParams, listObj []T) (*types.PageResponse, apperror.ErrorInfo)) (*types.PageResponse, apperror.ErrorInfo) {
	count, err := store.DB(ctx).GetCount(new(T), listParam)
	if err != nil {
		return nil, apperror.NewError(http.StatusBadRequest, "get count failed", err)
	}

	listObj := make([]T, 0)
	db := queryFn(store.GormDB(ctx), listParam)
	db = store.GenerateDBForQuery(db, listParam, &listObj)
	db = store.OptFuzzySearchDB(db, nil, listParam.FuzzySearchValue)
	db = store.OptPageDB(db, listParam)

	if err := db.Find(&listObj).Error; err != nil {
		return nil, apperror.NewError(http.StatusBadRequest, "load models failed", err)
	}

	return transFn(int(count), listParam, listObj)
}

func DefaultTransListObjToResp[T any, R any](count int, listParam *types.ListParams, listObj []T) (*types.PageResponse, apperror.ErrorInfo) {
	items := make([]interface{}, 0)

	for _, v := range listObj {
		r := new(R)

		err := utils.Copy(&r, &v)
		if err != nil {
			return nil, apperror.NewError(http.StatusBadRequest, "failed to copy", err)
		}

		items = append(items, r)
	}

	return types.NewPageResponse(count, listParam.PageInfo, items), nil
}
