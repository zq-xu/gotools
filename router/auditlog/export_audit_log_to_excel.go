package auditlog

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"zq-xu/helper/log"
	"zq-xu/helper/restapi"
	"zq-xu/helper/router/auth"
	"zq-xu/helper/utils"
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
	CreatedAt  utils.UnixTime
}

func init() {
	excelTittleList = make([]string, 0)

	t := reflect.TypeOf(ExcelRowData{})
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		excelTittleList = append(excelTittleList, t.Field(i).Name)
	}
}

func (erd *ExcelRowData) User(u auth.User) {
	erd.UserID = strconv.FormatInt(u.ID, 10)
	erd.UserName = u.Name
}

func ExportAuditLog(ctx *gin.Context) {
	listObj := make([]ModelAuditLog, 0)
	data := make([]interface{}, 0)
	conf := &restapi.ListConf{
		ModelObj:              &ModelAuditLog{},
		ModelObjList:          &listObj,
		FuzzySearchColumnList: []string{"url"},
		TransObjToRespFunc: func() interface{} {
			data = generateAuditLogListExcelData(listObj, data)
			return data
		},
		LoadAssociationsDBFunc: LoadAssociationsDB,
		ResponseWriteFunc:      func(ctx *gin.Context) { writeAuditLogToExcel(ctx, data) },
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func generateAuditLogListExcelData(objList []ModelAuditLog, data []interface{}) []interface{} {
	for _, v := range objList {
		r := ExcelRowData{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy audit log obj to response. %v", err)
			return nil
		}

		data = append(data, r)
	}

	return data
}

func writeAuditLogToExcel(ctx *gin.Context, data []interface{}) {
	ex := utils.NewExcelExport("")
	ex.WriteExcelByStruct("", excelTittleList, data)
	_ = ex.ExportExcelToGin("auditlog", ctx)
}
