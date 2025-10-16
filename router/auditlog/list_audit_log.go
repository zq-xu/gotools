package auditlog

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zq-xu/gotools"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/types"
)

type ResponseOfAuditLog struct {
	types.ModelResponse `json:",inline"`

	// User auth.ResponseOfUserInfo `json:"user"`

	ClientIP   string `json:"client_ip"`
	Url        string `json:"url"`
	Method     string `json:"method"`
	Body       string `json:"body"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func ListAuditLog(ctx *gin.Context) {
	reqParams, ei := gotools.GetListParams(ctx)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	resp, ei := listAuditLog(ctx, reqParams)
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func listAuditLog(ctx context.Context, listParam *types.ListParams) (*gotools.PageResponse, gotools.ErrorInfo) {
	listObj := make([]ModelAuditLog, 0)
	count, err := store.DB(ctx).ListWithCount(listParam, &ModelAuditLog{}, listObj)
	if err != nil {
		return nil, gotools.NewError(http.StatusBadRequest, "load models failed", err)
	}

	resp, ei := transObjToResp(listObj)
	if ei != nil {
		return nil, ei
	}

	return gotools.NewPageResponse(int(count), listParam.PageInfo, resp), nil
}

func transObjToResp(listObj []ModelAuditLog) ([]interface{}, gotools.ErrorInfo) {
	items := make([]interface{}, 0)

	for _, v := range listObj {
		r := ResponseOfAuditLog{}

		err := gotools.Copy(&r, &v)
		if err != nil {
			return nil, gotools.NewError(http.StatusBadRequest, "failed to copy", err)
		}

		items = append(items, r)
	}

	return items, nil
}
