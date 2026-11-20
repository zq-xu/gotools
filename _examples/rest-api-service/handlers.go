package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/apix"
	"github.com/zq-xu/gotools/apix/getx"
	"github.com/zq-xu/gotools/cryptox"
	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/storex"
	"github.com/zq-xu/gotools/storex/account"
	"github.com/zq-xu/gotools/storex/database"
	"github.com/zq-xu/gotools/utilsx"
)

type CreateAccountParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Roles    string `json:"roles"`
}

func createAccount(ctx *gin.Context) {
	p := new(CreateAccountParams)
	if err := ctx.ShouldBindJSON(p); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	cipherPwd, err := cryptox.Encrypt([]byte(p.Password))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	obj := &account.Account{
		Model:    storex.GenerateModel(),
		Username: p.Username,
		Password: cipherPwd,
		Name:     p.Name,
		Roles:    p.Roles,
	}

	ei := storex.DB(ctx).DoDBTransaction(func(db database.Database) errorx.ErrorInfo {
		ei := db.EnsureNotExistByField(obj, "username", p.Username)
		if ei != nil {
			return ei
		}

		if err := db.Create(obj); err != nil {
			return storex.NewErrorInfoForCreateError(err)
		}
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.StatusCode(), ei.ErrorMessage())
		return
	}

	ctx.JSON(http.StatusCreated, obj)
}

func listAccounts(ctx *gin.Context) {
	apix.ListHandler[account.Account](
		ctx,
		apix.DefaultListObjWithCount[account.Account],
		apix.DefaultTransListObjToResp[account.Account, account.Account],
	)
}

func getAccount(ctx *gin.Context) {
	getx.GetHandler[account.Account, account.Account](
		ctx,
		getx.NewDefaultConfig[account.Account, account.Account](),
	)
}

type UpdateAccountParams struct {
	Name  *string `json:"name"`
	Roles *string `json:"roles"`
}

func updateAccount(ctx *gin.Context) {
	apix.UpdateHandler[account.Account, UpdateAccountParams](
		ctx,
		func(db database.Database, obj *account.Account, p *UpdateAccountParams) errorx.ErrorInfo {
			utilsx.OptStringPtr(&obj.Name, p.Name)
			utilsx.OptStringPtr(&obj.Roles, p.Roles)
			return nil
		},
	)
}

func deleteAccount(ctx *gin.Context) {
	apix.DeleteHandler[account.Account](
		ctx,
		func(db database.Database, obj *account.Account) errorx.ErrorInfo {
			// 例如内置管理员不允许删除：这里返回非 nil 即可中断。
			return nil
		},
	)
}

// seedAdmin 给空库写入第一个可登录账号，由 main 的 -seed 参数触发。
func seedAdmin(ctx context.Context, username, rawPassword, roles string) error {
	obj := &account.Account{
		Model:    storex.GenerateModel(),
		Username: username,
		Password: mustEncrypt(rawPassword),
		Name:     username,
		Roles:    roles,
	}

	return storex.DB(ctx).DoDBTransaction(func(db database.Database) errorx.ErrorInfo {
		ei := db.EnsureNotExistByField(obj, "username", username)
		if ei != nil {
			return ei
		}

		if err := db.Create(obj); err != nil {
			return storex.NewErrorInfoForCreateError(err)
		}
		return nil
	})
}

func mustEncrypt(s string) string {
	cipher, err := cryptox.Encrypt([]byte(s))
	if err != nil {
		panic(err)
	}
	return cipher
}
