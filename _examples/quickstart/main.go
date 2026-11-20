package main

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/apix"
	"github.com/zq-xu/gotools/configx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/storex/account" // 导入即注册 accounts 表与 crypto 配置段
	"github.com/zq-xu/gotools/utilsx"
)

func main() {
	if err := configx.Setup("config.yaml"); err != nil {
		panic(err)
	}

	// 配置里没写 host/port 时补 0.0.0.0:8080
	routerx.RouterConfig.OptDefalt()

	routerx.RegisterGroup(
		routerx.NewGroup("/api/v1").AddAPI(
			&routerx.API{Method: "GET", Path: "/accounts", Handler: listAccounts},
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	utilsx.Notify(cancel) // 把 SIGINT / SIGTERM 接到 ctx 上

	if err := routerx.StartRouter(ctx, routerx.NewRouter()); err != nil {
		panic(err)
	}
}

func listAccounts(ctx *gin.Context) {
	apix.ListHandler[account.Account](
		ctx,
		apix.DefaultListObjWithCount[account.Account],
		apix.DefaultTransListObjToResp[account.Account, account.Account],
	)
}
