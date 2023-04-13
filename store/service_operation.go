package store

import (
	"errors"

	"gorm.io/gorm"

	"zq-xu/helper/restapi/response"
)

func DoDBTransaction(db *gorm.DB, fns ...func(db *gorm.DB) *response.ErrorInfo) *response.ErrorInfo {
	db = db.Begin()

	for _, fn := range fns {
		ei := fn(db)
		if ei != nil {
			db.Rollback()
			return ei
		}
	}

	err := db.Commit().Error
	if err != nil {
		return response.NewStorageError(response.TransactionCommitErrorCode, err.Error())
	}
	return nil
}

func EnsureExistByID(db *gorm.DB, obj interface{}, id string) *response.ErrorInfo {
	err := GetByID(db, obj, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCommonError(response.NotFoundErrorCode)
		}

		return response.NewStorageError(response.StorageErrorCode, err)
	}

	return nil
}

func EnsureNotExistByName(db *gorm.DB, obj interface{}, name string) *response.ErrorInfo {
	err := GetByName(db, obj, name)
	if err == nil {
		return response.NewCommonError(response.AlreadyExistsErrCode)

	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return response.NewStorageError(response.StorageErrorCode, err)
}
