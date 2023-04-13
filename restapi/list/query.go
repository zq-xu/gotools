package list

import "github.com/gin-gonic/gin"

// Request Query
const (
	QueryValue = "queryValue"
	BaseInfo   = "baseInfo"
)

var listQueryList = []string{QueryValue, BaseInfo}

type Queries map[string]string

//func GetQueries(ctx *gin.Context) map[string]string {
//	r := make(map[string]string, 0)
//	for _, v := range listQueryList {
//		s := ctx.Query(v)
//		if s != "" {
//			r[v] = s
//		}
//	}
//
//	return r
//}

func GetQueries(c *gin.Context) map[string]string {
	query := c.Request.URL.Query()
	var queryMap = make(map[string]string, len(query))
	for k := range query {
		queryMap[k] = c.Query(k)
	}
	return queryMap
}
