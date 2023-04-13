package list

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"zq-xu/helper/restapi/response"
)

const (
	defaultPageSize = 10

	PageNumParam  = "pageNum"
	PageSizeParam = "pageSize"
)

type PageResponse struct {
	PageInfo `json:",inline"`

	Count int         `json:"count"`
	Items interface{} `json:"items"`
}

type PageInfo struct {
	PageNum   int `json:"page_num"`
	PageSize  int `json:"page_size"`
	PageCount int `json:"page_total"`
}

func GetPageInfo(ctx *gin.Context) (*PageInfo, *response.ErrorInfo) {
	var err error
	pi := &PageInfo{}

	numStr := ctx.Query(PageNumParam)
	sizeStr := ctx.Query(PageSizeParam)

	if numStr != "" {
		pi.PageNum, err = strconv.Atoi(numStr)
		if err != nil {
			return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("PageNum is invalid. %v", err))
		}
	}

	if sizeStr != "" {
		pi.PageSize, err = strconv.Atoi(sizeStr)
		if err != nil {
			return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("PageSize is invalid. %v", err))
		}
	}

	pi.Revise()
	return pi, nil
}

func NewPageResponse(count int, pi *PageInfo, items interface{}) *PageResponse {
	pr := &PageResponse{
		PageInfo: *pi,
		Count:    count,
		Items:    items,
	}

	if pr.PageSize != 0 {
		pr.PageCount = pr.Count / pr.PageSize
	}

	if pr.PageCount*pr.PageSize < pr.Count {
		pr.PageCount++
	}

	return pr
}

func (p *PageInfo) Revise() {
	if p.PageSize == 0 {
		p.PageSize = defaultPageSize
	}

	if p.PageNum <= 0 {
		p.PageNum = 1
	}
}
