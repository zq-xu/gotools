package errorx

import (
	"fmt"
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

func Errorf(err error, status int, format string, msg ...string) ErrorInfo {
	return &errorInfo{status, fmt.Sprintf(format, msg), err}
}

func NewError(status int, message string, err error) ErrorInfo {
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
