package auditlog

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/routerx/auditlog/export"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/typesx"
	"github.com/zq-xu/gotools/utilsx"
)

var excelTittleList []string

type ExcelRowData struct {
	UserID     string
	UserName   string
	ClientIP   string
	Url        string
	Method     string
	Body       string
	Message    string
	StatusCode int
	CreatedAt  utilsx.UnixTime
}

// func (erd *ExcelRowData) User(u auth.User) {
// 	erd.UserID = strconv.FormatInt(u.ID, 10)
// 	erd.UserName = u.Name
// }

func init() {
	excelTittleList = make([]string, 0)

	t := reflect.TypeOf(ExcelRowData{})
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		excelTittleList = append(excelTittleList, t.Field(i).Name)
	}
}

func ExportAuditLog(ctx *gin.Context) {
	reqParams, ei := typesx.GetListParams(ctx)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	data, ei := exportAuditLog(ctx, reqParams)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	writeResponse(ctx, data)
}

func exportAuditLog(ctx context.Context, listParam *typesx.ListParams) ([]interface{}, errorx.ErrorInfo) {
	listObj := make([]ModelAuditLog, 0)
	err := storex.DB(ctx).List(listParam, listObj)
	if err != nil {
		return nil, errorx.NewError(http.StatusBadRequest, "load models failed", err)
	}

	return getExcelData(listObj)
}

func getExcelData(listObj []ModelAuditLog) ([]interface{}, errorx.ErrorInfo) {
	data := make([]interface{}, 0)

	for _, v := range listObj {
		r := ExcelRowData{}

		err := utilsx.Copy(&r, &v)
		if err != nil {
			return nil, errorx.NewError(http.StatusBadRequest, "failed to copy", err)
		}

		data = append(data, r)
	}

	return data, nil
}

func writeResponse(ctx *gin.Context, data []interface{}) {
	ex := export.NewExcelExport("")
	ex.WriteExcelByStruct("", excelTittleList, data)
	_ = ex.ExportExcelToGin("auditlog", ctx)
}
