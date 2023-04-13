package restapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zq-xu/helper/restapi/response"
	"zq-xu/helper/store"
	"zq-xu/warehouse-admin/internal/webserver/types"
)

type DetailConf struct {
	ModelObj               interface{}
	RespObj                interface{}
	TransObjToRespFunc     func() interface{}
	LoadAssociationsDBFunc func(db *gorm.DB) *gorm.DB
}

func GetDetail(ctx *gin.Context, conf *DetailConf) {
	id, ei := getID(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = getDetail(ctx, id, conf)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusOK, conf.RespObj)
}

func getID(ctx *gin.Context) (int64, *response.ErrorInfo) {
	id := ctx.Param(types.IDParam)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, response.NewCommonError(response.InvalidParametersErrorCode,
			fmt.Sprintf("invalid id. %v", err))
	}

	return idInt, nil
}

func getDetail(ctx context.Context, id int64, conf *DetailConf) *response.ErrorInfo {
	db := store.DB(ctx)

	err := conf.LoadAssociationsDBFunc(db).
		Preload(clause.Associations).
		Where("id = ?", id).
		First(conf.ModelObj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCommonError(response.NotFoundErrorCode)
		}
		return response.NewStorageError(response.StorageErrorCode, err)
	}

	conf.TransObjToRespFunc()
	return nil
}
