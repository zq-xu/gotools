package apix

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/logx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/storex/database"
)

func DeleteHandler[T any](ctx *gin.Context, validateFn func(database.Database, *T) errorx.ErrorInfo) {
	id := routerx.GetID(ctx)

	ei := delete(ctx, id, validateFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), errorx.NewResponseBody(ei))
		return
	}

	tType := reflect.TypeOf((*T)(nil)).Elem()
	logx.Logger.Infof("Succeed to delete obj %s/%s", tType, id)

	ctx.JSON(http.StatusNoContent, struct{}{})
}

func delete[T any](ctx *gin.Context, id string, validateFn func(database.Database, *T) errorx.ErrorInfo) errorx.ErrorInfo {
	return storex.DB(ctx).DoDBTransaction(func(db database.Database) errorx.ErrorInfo {
		obj := new(T)
		err := db.GetAssociations(obj, id)
		if storex.IsNotFoundError(err) {
			return nil
		}
		if err != nil {
			return errorx.NewError(http.StatusInternalServerError, "unexpected error during get detail", err)
		}

		if validateFn != nil {
			ei := validateFn(db, obj)
			if ei != nil {
				return ei
			}
		}

		err = db.DeleteAssociations(obj, id)
		if err != nil {
			return errorx.NewError(http.StatusBadRequest, "delete failed", err)
		}

		return nil
	})
}
