package gormkit

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/zq-xu/gotools/errorx"
)

func NewErrorInfoForGetError(err error) errorx.ErrorInfo {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.NewError(http.StatusBadRequest, "not found", err)
		}

		return errorx.NewError(http.StatusInternalServerError, "unexpected error during get detail", err)
	}

	return nil
}

func NewErrorInfoForListError(err error) errorx.ErrorInfo {
	if err != nil {
		return errorx.NewError(http.StatusBadRequest, "unexpected error during list", err)
	}
	return nil
}

func IsNotFoundError(err error) bool {
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}

	return false
}

func NewErrorInfoForUpdateError(err error) errorx.ErrorInfo {
	if err != nil {
		return errorx.NewError(http.StatusInternalServerError, "unexpected error during update", err)
	}

	return nil
}

func NewErrorInfoForCreateError(err error) errorx.ErrorInfo {
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errorx.NewError(http.StatusBadRequest, "update failed", err)
		}

		return errorx.NewError(http.StatusInternalServerError, "unexpected error during creation", err)
	}

	return nil
}

func NewErrorInfoForDeleteError(err error) errorx.ErrorInfo {
	if err != nil {
		return errorx.NewError(http.StatusInternalServerError, "unexpected error during delete", err)
	}

	return nil
}
