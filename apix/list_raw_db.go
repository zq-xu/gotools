package apix

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/typesx"
	"github.com/zq-xu/gotools/utilsx"
)

func ListByRawGormHandler[T any](ctx *gin.Context,
	queryFn func(db *gorm.DB, listParam *typesx.ListParams) *gorm.DB,
	transFn func(count int, listParam *typesx.ListParams, listObj []T) (*typesx.PageResponse, errorx.ErrorInfo)) {
	listParam, ei := typesx.GetListParams(ctx)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), errorx.NewResponseBody(ei))
		return
	}

	resp, ei := listByRawGorm(ctx, listParam, queryFn, transFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), errorx.NewResponseBody(ei))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func listByRawGorm[T any](ctx context.Context, listParam *typesx.ListParams,
	queryFn func(db *gorm.DB, listParam *typesx.ListParams) *gorm.DB,
	transFn func(count int, listParam *typesx.ListParams, listObj []T) (*typesx.PageResponse, errorx.ErrorInfo)) (*typesx.PageResponse, errorx.ErrorInfo) {
	listObj := make([]T, 0)

	// 计数与列表使用同一套条件：queryFn 附加的 WHERE + 查询条件 + 模糊搜索，
	// 分页（Limit/Offset/Order）除外，保证 total 与列表一致。
	countDB := queryFn(storex.GormDB(ctx), listParam)
	countDB = storex.GenerateDBForQuery(countDB, listParam, &listObj)
	countDB = storex.OptFuzzySearchDB(countDB, nil, listParam.FuzzySearchValue)

	var count int64
	if err := countDB.Session(&gorm.Session{}).Model(new(T)).Count(&count).Error; err != nil {
		return nil, errorx.NewError(http.StatusBadRequest, "get count failed", err)
	}

	db := queryFn(storex.GormDB(ctx), listParam)
	db = storex.GenerateDBForQuery(db, listParam, &listObj)
	db = storex.OptFuzzySearchDB(db, nil, listParam.FuzzySearchValue)
	db = storex.OptPageDB(db, listParam)

	if err := db.Find(&listObj).Error; err != nil {
		return nil, errorx.NewError(http.StatusBadRequest, "load models failed", err)
	}

	return transFn(int(count), listParam, listObj)
}

func DefaultTransListObjToResp[T any, R any](count int, listParam *typesx.ListParams, listObj []T) (*typesx.PageResponse, errorx.ErrorInfo) {
	items := make([]interface{}, 0)

	for _, v := range listObj {
		r := new(R)

		err := utilsx.Copy(&r, &v)
		if err != nil {
			return nil, errorx.NewError(http.StatusBadRequest, "failed to copy", err)
		}

		items = append(items, r)
	}

	return typesx.NewPageResponse(count, listParam.PageInfo, items), nil
}
