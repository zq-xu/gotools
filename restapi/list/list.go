package list

import (
	"github.com/gin-gonic/gin"

	"zq-xu/helper/restapi/response"
)

const (
	SortByQuery = "sort_by"
)

type Params struct {
	PageInfo  *PageInfo
	Queries   Queries
	SortQuery string
}

func GetListParams(ctx *gin.Context) (*Params, *response.ErrorInfo) {
	pi, ei := GetPageInfo(ctx)
	if ei != nil {
		return nil, ei
	}

	return &Params{
		PageInfo:  pi,
		Queries:   GetQueries(ctx),
		SortQuery: ctx.Query(SortByQuery),
	}, nil
}
