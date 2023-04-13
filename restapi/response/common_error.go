package response

import "net/http"

const (
	InvalidParametersErrorCode ErrorSectionCode = iota
	AlreadyExistsErrCode
	NotFoundErrorCode
	GenerateModelErrorCode
	GetFormFileErrorCode
	SaveFileErrorCode
	CheckDirErrorCode
	UploadFileToS3ErrorCode
	ResizeFileErrorCode
	InvalidImageFormatErrorCode
)

var commonErrorList = []ServiceErrorInfo{
	{ErrorSectionCode: InvalidParametersErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "Invalid parameters!"}},
	{ErrorSectionCode: AlreadyExistsErrCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusConflict, ErrorMessageFmt: "The object already exists!"}},
	{ErrorSectionCode: NotFoundErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "The object is not found!"}},
	{ErrorSectionCode: GenerateModelErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Generate model error!"}},
	{ErrorSectionCode: GetFormFileErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "Get form file error!"}},
	{ErrorSectionCode: SaveFileErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Save file error!"}},
	{ErrorSectionCode: CheckDirErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Check dir error!"}},
	{ErrorSectionCode: UploadFileToS3ErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Upload file error!"}},
	{ErrorSectionCode: ResizeFileErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "resize file error!"}},
	{ErrorSectionCode: InvalidImageFormatErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "invalid file format error! support jpg / jpeg / png."}},
}

func init() {
	RegisterErrorInfo(CommonErrorSectionCode, commonErrorList...)
}
