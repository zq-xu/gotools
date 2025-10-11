package apperror

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type ErrorInfo interface {
	StatusCode() int
	ErrorMessage() string
	Error() string
}

type errorInfo struct {
	status  int
	message string
	err     error
}

func Errorf(err error, status int, format string, msg ...string) *errorInfo {
	return &errorInfo{status, fmt.Sprintf(format, msg), err}
}

func NewError(status int, message string, err error) *errorInfo {
	return &errorInfo{status, message, err}
}

func (ei *errorInfo) StatusCode() int {
	return ei.status
}

func (ei *errorInfo) Error() string {
	if ei.err == nil {
		return ei.message
	}

	return fmt.Sprintf("%s. %s", ei.message, ei.err.Error())
}

func (ei *errorInfo) ErrorMessage() string {
	return ei.Error()
}

func NewErrorInfoForDBGetError(err error) *errorInfo {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewError(http.StatusBadRequest, "not found", err)
		}

		return NewError(http.StatusInternalServerError, "unexpected error during get detail", err)
	}

	return nil
}

func NewErrorInfoForUpdateError(err error) *errorInfo {
	if err != nil {
		return NewError(http.StatusInternalServerError, "unexpected error during update", err)
	}

	return nil
}

func NewErrorInfoForCreateError(err error) *errorInfo {
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return NewError(http.StatusBadRequest, "update failed", err)
		}

		return NewError(http.StatusInternalServerError, "unexpected error during creation", err)
	}

	return nil
}

func NewErrorInfoForDeleteError(err error) *errorInfo {
	if err != nil {
		return NewError(http.StatusInternalServerError, "unexpected error during delete", err)
	}

	return nil
}
