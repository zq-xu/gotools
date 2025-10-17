package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/zq-xu/gotools/storex"
)

const (
	ModelAuditLogTableName = "audit_logs"
)

type ModelAuditLog struct {
	storex.Model `json:",inline"`

	UserID    int64
	UserName  string
	UserAlias string

	ClientIP   string
	Url        string
	Method     string
	Body       string `gorm:"type:text"`
	Message    string
	StatusCode int
}

func init() {
	storex.RegisterTable(&ModelAuditLog{})
}

func (al *ModelAuditLog) TableName() string {
	return ModelAuditLogTableName
}

func NewModelAuditLog(ctx *gin.Context, reqBody []byte, msg []byte) *ModelAuditLog {
	// id := ctx.GetString(auth.AuthUserIDToken)
	// idInt64, _ := strconv.ParseInt(id, 10, 64)

	// name := ctx.GetString(auth.AuthUserNameToken)
	// alias := ctx.GetString(auth.AuthUserUsernameToken)

	return &ModelAuditLog{
		Model: storex.GenerateModel(),
		// UserID:     idInt64,
		// UserName:   name,
		// UserAlias:  alias,
		Method:     ctx.Request.Method,
		Url:        ctx.Request.URL.Path,
		ClientIP:   ctx.ClientIP(),
		StatusCode: ctx.Writer.Status(),
		Body:       string(reqBody),
		Message:    string(msg),
	}
}
