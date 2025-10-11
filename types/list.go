package types

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"

	"zq-xu/gotools/apperror"
)

const (
	defaultPageSize = 10

	PageNumParam     = "pageNum"
	PageSizeParam    = "pageSize"
	FuzzySearchParam = "fuzzySearch"
)

var skipQueryKeyList = []string{PageNumParam, PageSizeParam, sortByQuery}

type ListParams struct {
	PageInfo *PageInfo
	Queries  Queries
	Sorter   *sorter

	FuzzySearchColumnList []string
	FuzzySearchValue      string
}

type Queries map[string]string

type PageResponse struct {
	PageInfo `json:",inline"`

	Count int         `json:"count"`
	Items interface{} `json:"items"`
}

type PageInfo struct {
	PageNum   int `json:"pageNum"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`
}

func GetListParams(ctx *gin.Context) (*ListParams, apperror.ErrorInfo) {
	pi, ei := getPageInfo(ctx)
	if ei != nil {
		return nil, ei
	}

	return &ListParams{
		PageInfo:         pi,
		Queries:          getQueries(ctx),
		Sorter:           NewSorter(ctx.DefaultQuery(sortByQuery, defaultSortStr)),
		FuzzySearchValue: ctx.Query(FuzzySearchParam),
	}, nil
}

func (l *ListParams) Validate(obj interface{}) apperror.ErrorInfo {
	return l.Sorter.validateSortQuery(obj)
}

func getQueries(c *gin.Context) map[string]string {
	query := c.Request.URL.Query()
	var queryMap = make(map[string]string, len(query))

	for k := range query {
		if slices.Contains(skipQueryKeyList, k) {
			continue
		}

		queryMap[k] = c.Query(k)
	}
	return queryMap
}

func getPageInfo(ctx *gin.Context) (*PageInfo, apperror.ErrorInfo) {
	var err error
	pi := &PageInfo{}

	numStr := ctx.Query(PageNumParam)
	sizeStr := ctx.Query(PageSizeParam)

	if numStr != "" {
		pi.PageNum, err = strconv.Atoi(numStr)
		if err != nil {
			return nil, apperror.NewError(http.StatusBadRequest, "PageNum is invalid", err)
		}
	}

	if sizeStr != "" {
		pi.PageSize, err = strconv.Atoi(sizeStr)
		if err != nil {
			return nil, apperror.NewError(http.StatusBadRequest, "PageSize is invalid", err)
		}
	}

	pi.Revise()
	return pi, nil
}

// NewPageResponse
func NewPageResponse(count int, pi *PageInfo, items []interface{}) *PageResponse {
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
