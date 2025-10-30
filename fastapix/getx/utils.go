package getx

import (
	"net/http"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/utilsx"
)

func DefaultTransObjToResp[T any, R any](obj *T) (*R, errorx.ErrorInfo) {
	resp := new(R)

	err := utilsx.Copy(resp, obj)
	if err != nil {
		return nil, errorx.NewError(http.StatusBadRequest, "failed to copy", err)
	}
	return resp, nil
}
