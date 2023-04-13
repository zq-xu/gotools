package response

import (
	"fmt"
	"net/http"

	"zq-xu/helper/log"
)

const ErrorBaseCode = "E"

const (
	// the service error begins with zero
	ServiceErrorSectionCode ErrorSectionCode = iota
	CommonErrorSectionCode
	GlobalErrorSectionCode
	NetworkErrorSectionCode
	StorageErrorSectionCode
	ResourceErrorSectionCode
)

const (
	InternalErrorCode ErrorSectionCode = iota
)

const (
	StorageErrorCode ErrorSectionCode = iota
	TransactionCommitErrorCode
)

var globalErrorList = []ServiceErrorInfo{
	{ErrorSectionCode: InternalErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Internal error! %v"}},
	{ErrorSectionCode: InvalidParametersErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "Invalid parameters! %v"}},
}

var storageErrorList = []ServiceErrorInfo{
	{ErrorSectionCode: StorageErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Failed to operate storage! %v"}},
	{ErrorSectionCode: TransactionCommitErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Failed to commit the transaction! %v"}},
}

var (
	errorSet = make(map[ErrorCode]ErrorInfo)

	UnknownError = &ErrorInfo{
		ErrorCode: "unknown",
		ErrorBaseInfo: ErrorBaseInfo{
			ErrorMessage: "unknown",
		},
	}
)

// fmt.Sprintf("%s%d%d%d", ErrorBaseCode, ErrType, ServiceType, ErrorSectionCode)
type ErrorCode string

type ErrorSectionCode int

type ErrorBaseInfo struct {
	ErrorMessage string `json:"errorMessage"`

	Status          int    `json:"-"`
	ErrorMessageFmt string `json:"-"`
}

type ErrorInfo struct {
	ErrorCode ErrorCode `json:"errorCode"`
	ErrorBaseInfo
}

type ServiceErrorInfo struct {
	ErrorSectionCode ErrorSectionCode
	ErrorBaseInfo
}

func NewErrorInfo(errCode ErrorCode, msg ...interface{}) *ErrorInfo {
	v, ok := errorSet[errCode]
	if !ok {
		return UnknownError
	}

	v.ErrorMessage = fmt.Sprintf(v.ErrorMessageFmt, msg...)

	return &v
}

func (ei *ErrorInfo) Error() string {
	return ei.ErrorMessage
}

func init() {
	RegisterErrorInfo(GlobalErrorSectionCode, globalErrorList...)
	RegisterErrorInfo(StorageErrorSectionCode, storageErrorList...)
}

func GenerateErrorCode(errType, code ErrorSectionCode) ErrorCode {
	return ErrorCode(fmt.Sprintf("%s%04d%04d", ErrorBaseCode, errType, code))
}

func RegisterErrorInfo(errType ErrorSectionCode, list ...ServiceErrorInfo) {
	for _, sei := range list {
		ec := GenerateErrorCode(errType, sei.ErrorSectionCode)
		_, ok := errorSet[ec]
		if ok {
			log.Logger.Warningf("error code %v has already exist!", ec)
			continue
		}

		errorSet[ec] = ErrorInfo{
			ErrorCode:     ec,
			ErrorBaseInfo: sei.ErrorBaseInfo,
		}
	}
}

func NewCommonError(errCode ErrorSectionCode, msg ...interface{}) *ErrorInfo {
	ec := GenerateErrorCode(CommonErrorSectionCode, errCode)
	return NewErrorInfo(ec, msg...)
}

func NewGlobalError(errCode ErrorSectionCode, msg ...interface{}) *ErrorInfo {
	ec := GenerateErrorCode(GlobalErrorSectionCode, errCode)
	return NewErrorInfo(ec, msg...)
}

func NewStorageError(errCode ErrorSectionCode, msg ...interface{}) *ErrorInfo {
	ec := GenerateErrorCode(StorageErrorSectionCode, errCode)
	return NewErrorInfo(ec, msg...)
}

func NewServiceError(errCode ErrorSectionCode, msg ...interface{}) *ErrorInfo {
	ec := GenerateErrorCode(ServiceErrorSectionCode, errCode)
	return NewErrorInfo(ec, msg...)
}
