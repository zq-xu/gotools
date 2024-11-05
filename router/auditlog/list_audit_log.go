package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/helper/log"
	"zq-xu/helper/restapi"
	"zq-xu/helper/router/auth"
	"zq-xu/helper/store"
)

type ResponseOfAuditLog struct {
	store.Model `json:",inline"`

	User auth.ResponseOfUserInfo `json:"user"`

	ClientIP   string `json:"client_ip"`
	Url        string `json:"url"`
	Method     string `json:"method"`
	Body       string `json:"body"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type ListResponseOfAuditLog []ResponseOfAuditLog

func ListAuditLog(ctx *gin.Context) {
	listObj := make([]ModelAuditLog, 0)

	conf := &restapi.ListConf{
		ModelObj:               &ModelAuditLog{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"url"},
		TransObjToRespFunc:     func() interface{} { return generateAuditLogListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB) *gorm.DB {
	return GenerateReadAuditLogDB(db)
}

func generateAuditLogListResponse(objList []ModelAuditLog) interface{} {
	items := make(ListResponseOfAuditLog, 0)

	for _, v := range objList {
		r := ResponseOfAuditLog{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy audit log obj to response. %v", err)
			return nil
		}

		items = append(items, r)
	}

	return items
}
