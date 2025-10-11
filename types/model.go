package types

import "zq-xu/gotools/utils"

type ModelResponse struct {
	ID        string         `json:"id"`
	CreatedAt utils.UnixTime `json:"createdAt"`
	UpdatedAt utils.UnixTime `json:"updatedAt"`
	Comment   string         `json:"comment"`
	Status    int64          `json:"status"`
}
