package fastapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/router"
	"github.com/zq-xu/gotools/store"
)

func UpdateByRawGormHandler[T any, P any](ctx *gin.Context,
	queryFn func(*gorm.DB, string) *gorm.DB,
	optFn func(*gorm.DB, *T, *P) gotools.ErrorInfo,
	afterUpdateFn func(*gorm.DB, *T, *P) gotools.ErrorInfo) {
	id := router.GetID(ctx)

	reqParams := new(P)
	err := ctx.ShouldBindJSON(reqParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid parameters. %s", err))
		return
	}

	ei := updateByRawGorm(ctx, id, reqParams, queryFn, optFn, afterUpdateFn)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func updateByRawGorm[T any, P any](ctx context.Context, id string, reqParams *P,
	queryFn func(*gorm.DB, string) *gorm.DB,
	optFn func(*gorm.DB, *T, *P) apperror.ErrorInfo,
	afterUpdateFn func(*gorm.DB, *T, *P) apperror.ErrorInfo) apperror.ErrorInfo {
	return store.DoGormDBTransaction(store.GormDB(ctx), func(db *gorm.DB) apperror.ErrorInfo {
		obj := new(T)
		err := queryFn(store.GormDB(ctx), id).First(obj).Error
		ei := store.NewErrorInfoForGetError(err)
		if ei != nil {
			return ei
		}

		ei = optFn(db, obj, reqParams)
		if ei != nil {
			return ei
		}

		err = db.Omit(clause.Associations).Save(obj).Error
		if err != nil {
			return apperror.NewError(http.StatusBadRequest, "update failed", err)
		}

		return afterUpdateFn(db, obj, reqParams)
	})
}
