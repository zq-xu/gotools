package gormkit

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/zq-xu/gotools/apperror"
)

func NewErrorInfoForGetError(err error) apperror.ErrorInfo {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NewError(http.StatusBadRequest, "not found", err)
		}

		return apperror.NewError(http.StatusInternalServerError, "unexpected error during get detail", err)
	}

	return nil
}

func IsNotFoundError(err error) bool {
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}

	return false
}

func NewErrorInfoForUpdateError(err error) apperror.ErrorInfo {
	if err != nil {
		return apperror.NewError(http.StatusInternalServerError, "unexpected error during update", err)
	}

	return nil
}

func NewErrorInfoForCreateError(err error) apperror.ErrorInfo {
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperror.NewError(http.StatusBadRequest, "update failed", err)
		}

		return apperror.NewError(http.StatusInternalServerError, "unexpected error during creation", err)
	}

	return nil
}

func NewErrorInfoForDeleteError(err error) apperror.ErrorInfo {
	if err != nil {
		return apperror.NewError(http.StatusInternalServerError, "unexpected error during delete", err)
	}

	return nil
}
