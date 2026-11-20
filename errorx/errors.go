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

// Errorf 组装一条带状态码的错误消息，args 用于填充 format 里的占位符。
func Errorf(err error, status int, format string, args ...any) ErrorInfo {
	return &errorInfo{status, fmt.Sprintf(format, args...), err}
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
