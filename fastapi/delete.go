package fastapi

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/router"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/store/database"
)

func DeleteHandler[T any](ctx *gin.Context, validateFn func(database.Database, *T) gotools.ErrorInfo) {
	id := router.GetID(ctx)

	ei := delete[T](ctx, id, validateFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	tType := reflect.TypeOf((*T)(nil)).Elem()
	gotools.Logger.Infof("Succeed to delete obj %s/%s", tType, id)

	ctx.JSON(http.StatusNoContent, struct{}{})
}

func delete[T any](ctx *gin.Context, id string, validateFn func(database.Database, *T) gotools.ErrorInfo) gotools.ErrorInfo {
	return store.DB(ctx).DoDBTransaction(func(db database.Database) gotools.ErrorInfo {
		obj := new(T)
		err := db.GetAssociations(obj, id)
		if store.IsNotFoundError(err) {
			return nil
		}
		if err != nil {
			return apperror.NewError(http.StatusInternalServerError, "unexpected error during get detail", err)
		}

		ei := validateFn(db, obj)
		if ei != nil {
			return ei
		}

		err = db.DeleteAssociations(obj, id)
		if err != nil {
			return gotools.NewError(http.StatusBadRequest, "delete failed", err)
		}

		return nil
	})
}
