package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/helper/restapi/response"
	"zq-xu/helper/store"
)

var (
	RoleSet = map[int]string{
		0: "user",
		1: "admin",
		2: "super",
	}
)

type ResponseOfUserInfo struct {
	Name     string `json:"name"`
	UserName string `json:"userName"`
	UserRole string `json:"userRole"`
	Status   int    `json:"status"`
}

func (r *ResponseOfUserInfo) Role(rl int) {
	r.UserRole = RoleSet[rl]
}

func (r *ResponseOfUserInfo) Alias(a string) {
	r.UserName = a
}

func GetUserInfo(ctx *gin.Context) {
	id := ctx.GetString(AuthUserIDToken)

	obj, ei := getUserModel(ctx, id)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	res, ei := generateUserResponse(obj)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func getUserModel(ctx context.Context, id string) (*User, *response.ErrorInfo) {
	db := store.DB(ctx)

	obj := &User{}
	err := db.Where("id = ?", id).First(obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCommonError(response.NotFoundErrorCode)
		}

		return nil, response.NewStorageError(response.StorageErrorCode, err)
	}
	return obj, nil
}

func generateUserResponse(obj *User) (*ResponseOfUserInfo, *response.ErrorInfo) {
	resp := &ResponseOfUserInfo{}

	err := copier.Copy(resp, obj)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return resp, nil
}
