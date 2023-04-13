package auditlog

import (
	"gorm.io/gorm"
	"strconv"

	"github.com/gin-gonic/gin"

	"zq-xu/helper/router/auth"
	"zq-xu/helper/store"
)

const (
	ModelAuditLogTableName = "audit_log"
)

type ModelAuditLog struct {
	store.Model

	// The ModelAuditLog belongs to the User
	UserID int64
	User   auth.User

	ClientIP   string
	Url        string
	Method     string
	Body       string
	Message    string
	StatusCode int
}

func init() {
	store.RegisterTable(&ModelAuditLog{})
}

func (al *ModelAuditLog) TableName() string {
	return ModelAuditLogTableName
}

func GenerateReadAuditLogDB(db *gorm.DB) *gorm.DB {
	return db.Preload("User")
}

func NewModelAuditLog(ctx *gin.Context, reqBody []byte, msg []byte) *ModelAuditLog {
	id := ctx.GetString(auth.AuthUserIDToken)
	idInt64, _ := strconv.ParseInt(id, 10, 64)

	return &ModelAuditLog{
		Model:      store.GenerateModel(),
		UserID:     idInt64,
		Method:     ctx.Request.Method,
		Url:        ctx.Request.URL.Path,
		ClientIP:   ctx.ClientIP(),
		StatusCode: ctx.Writer.Status(),
		Body:       string(reqBody),
		Message:    string(msg),
	}
}
