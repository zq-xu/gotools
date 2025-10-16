package auditlog

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/router/auditlog/export"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/types"
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
	CreatedAt  gotools.UnixTime
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
	reqParams, ei := gotools.GetListParams(ctx)
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

func exportAuditLog(ctx context.Context, listParam *types.ListParams) ([]interface{}, gotools.ErrorInfo) {
	listObj := make([]ModelAuditLog, 0)
	err := store.DB(ctx).List(listParam, listObj)
	if err != nil {
		return nil, gotools.NewError(http.StatusBadRequest, "load models failed", err)
	}

	return getExcelData(listObj)
}

func getExcelData(listObj []ModelAuditLog) ([]interface{}, gotools.ErrorInfo) {
	data := make([]interface{}, 0)

	for _, v := range listObj {
		r := ExcelRowData{}

		err := gotools.Copy(&r, &v)
		if err != nil {
			return nil, gotools.NewError(http.StatusBadRequest, "failed to copy", err)
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
