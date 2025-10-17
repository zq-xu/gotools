package routerx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	IDParam = "id"
)

// GetID
func GetID(ctx *gin.Context) string {
	return ctx.Param(IDParam)
}

// GetIDInt64FromToken
func GetIDInt64FromToken(ctx *gin.Context) (int64, error) {
	str := GetUserInfoFromToken(ctx).ID
	return strconv.ParseInt(str, 10, 64)
}
