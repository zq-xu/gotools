package fastapix

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/storex/database"
)

func UpdateHandler[T any, P any](ctx *gin.Context, optFn func(db database.Database, obj *T, params *P) errorx.ErrorInfo) {
	id := routerx.GetID(ctx)

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

func update[T any, P any](ctx *gin.Context, id string, reqParams *P, optFn func(db database.Database, obj *T, params *P) errorx.ErrorInfo) errorx.ErrorInfo {
	return storex.DB(ctx).DoDBTransaction(func(db database.Database) errorx.ErrorInfo {
		obj := new(T)
		err := db.Get(obj, id)
		ei := storex.NewErrorInfoForGetError(err)
		if ei != nil {
			return ei
		}

		ei = optFn(db, obj, reqParams)
		if ei != nil {
			return ei
		}

		if err := db.Update(obj); err != nil {
			return errorx.NewError(http.StatusBadRequest, "update failed", err)
		}
		return nil
	})
}
