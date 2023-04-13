package restapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/helper/log"
	"zq-xu/helper/restapi/response"
	"zq-xu/helper/store"
)

func Delete(ctx *gin.Context, obj interface{}) {
	id := ctx.Param(IDParam)

	ei := store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		err := store.GetByID(db, obj, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			return response.NewStorageError(response.StorageErrorCode, err)
		}

		err = store.DeleteAssociationsByID(db, obj, id)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	log.Logger.Infof("Succeed to delete obj %+v/%d", obj, id)
	ctx.JSON(http.StatusNoContent, struct{}{})
}
