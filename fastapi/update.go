package fastapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/router"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/store/database"
)

func UpdateHandler[T any, P any](ctx *gin.Context, optFn func(db database.Database, obj *T, params *P) gotools.ErrorInfo) {
	id := router.GetID(ctx)

	reqParams := new(P)
	err := ctx.ShouldBindJSON(reqParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid parameters. %s", err))
		return
	}

	ei := update(ctx, id, reqParams, optFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func update[T any, P any](ctx *gin.Context, id string, reqParams *P, optFn func(db database.Database, obj *T, params *P) gotools.ErrorInfo) gotools.ErrorInfo {
	return store.DB(ctx).DoDBTransaction(func(db database.Database) gotools.ErrorInfo {
		obj := new(T)
		err := db.Get(obj, id)
		ei := store.NewErrorInfoForGetError(err)
		if ei != nil {
			return ei
		}

		ei = optFn(db, obj, reqParams)
		if ei != nil {
			return ei
		}

		if err := db.Update(obj); err != nil {
			return gotools.NewError(http.StatusBadRequest, "update failed", err)
		}
		return nil
	})
}
