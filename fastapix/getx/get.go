package getx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/logx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/storex/database"
)

type Config[T any, R any] struct {
	QueryFn func(db database.Database, id string) (*T, errorx.ErrorInfo)
	TransFn func(obj *T) (*R, errorx.ErrorInfo)
}

func (cfg *Config[T, R]) MustValidate(key string) {
	if cfg.QueryFn == nil {
		logx.Logger.Fatalf("%s init failed. empty QueryFn", key)
		return
	}

	if cfg.TransFn == nil {
		logx.Logger.Fatalf("%s init failed. empty TransFn", key)
		return
	}
}

func NewDefaultConfig[T any, R any]() *Config[T, R] {
	return &Config[T, R]{
		QueryFn: DefaultQueryObj[T],
		TransFn: DefaultTransObjToResp[T, R],
	}
}

func DefaultQueryObj[T any](db database.Database, id string) (*T, errorx.ErrorInfo) {
	obj := new(T)
	ei := storex.NewErrorInfoForGetError(db.GetAssociations(obj, id))
	if ei != nil {
		return nil, ei
	}
	return obj, nil
}

func GetHandler[T any, R any](ctx *gin.Context, cfg *Config[T, R]) {
	id := routerx.GetID(ctx)

	obj, ei := cfg.QueryFn(storex.DB(ctx), id)
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
