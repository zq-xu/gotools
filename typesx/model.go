package typesx

import "github.com/zq-xu/gotools/utilsx"

type ModelResponse struct {
	ID        string          `json:"id"`
	CreatedAt utilsx.UnixTime `json:"createdAt"`
	UpdatedAt utilsx.UnixTime `json:"updatedAt"`
	Comment   string          `json:"comment"`
	Status    int64           `json:"status"`
}
