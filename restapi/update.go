package restapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/helper/log"
	"zq-xu/helper/restapi/response"
	"zq-xu/helper/store"
)

type UpdateConf struct {
	UpdateReq        interface{}
	ModelObj         interface{}
	OptModelFunc     func(db *gorm.DB) *response.ErrorInfo
	DealAssociations func(db *gorm.DB) *response.ErrorInfo
}

func Update(ctx *gin.Context, conf *UpdateConf) {
	id := ctx.Param(IDParam)
	err := ctx.ShouldBindJSON(conf.UpdateReq)
	if err != nil {
		ei := response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
		ctx.JSON(ei.Status, ei)
		return
	}

	ei := store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		err := store.GetByID(db, conf.ModelObj, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewCommonError(response.NotFoundErrorCode)
			}

			return response.NewStorageError(response.StorageErrorCode, err)
		}

		ei := conf.OptModelFunc(db)
		if ei != nil {
			return ei
		}

		err = store.Update(db, conf.ModelObj)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		if conf.DealAssociations != nil {
			ei = conf.DealAssociations(db)
			if ei != nil {
				return ei
			}
		}

		log.Logger.Infof("Succeed to update obj %+v", conf.ModelObj)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}
