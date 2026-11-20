package getx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/logx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/storex"
)

type RawGormConfig[T any, R any] struct {
	ValidateFn validateFn
	QueryFn    func(db *gorm.DB, id string) (*T, errorx.ErrorInfo)
	TransFn    func(obj *T) (*R, errorx.ErrorInfo)
}

func (cfg *RawGormConfig[T, R]) MustValidate(key string) {
	if cfg.QueryFn == nil {
		logx.Logger.Fatalf("%s init failed. empty QueryFn", key)
		return
	}

	if cfg.TransFn == nil {
		logx.Logger.Fatalf("%s init failed. empty TransFn", key)
		return
	}
}

func NewDefaultRawGormConfig[T any, R any]() *RawGormConfig[T, R] {
	return &RawGormConfig[T, R]{
		QueryFn: DefaultQueryObjByRawGorm[T],
		TransFn: DefaultTransObjToResp[T, R],
	}
}

func DefaultQueryObjByRawGorm[T any](db *gorm.DB, id string) (*T, errorx.ErrorInfo) {
	obj := new(T)
	ei := storex.NewErrorInfoForGetError(db.Where("id = ?", id).First(obj).Error)
	if ei != nil {
		return nil, ei
	}
	return obj, nil
}

func GetByRawGormHandler[T any, R any](ctx *gin.Context,
	cfg *RawGormConfig[T, R]) {
	id := routerx.GetID(ctx)

	if cfg.ValidateFn != nil {
		ei := cfg.ValidateFn(ctx)
		if ei != nil {
			ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
			return
		}
	}

	obj, ei := cfg.QueryFn(storex.GormDB(ctx), id)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	resp, ei := cfg.TransFn(obj)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
