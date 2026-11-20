package main

import (
	"context"
	"flag"
	"os"

	"github.com/zq-xu/gotools/configx"
	"github.com/zq-xu/gotools/logx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/routerx/auditlog"
	"github.com/zq-xu/gotools/storex/account"
	"github.com/zq-xu/gotools/utilsx"
)

func main() {
	seedOnly := flag.Bool("seed", false, "创建初始管理员账号 admin/admin123 后退出")
	flag.Parse()

	if err := configx.Setup("config.yaml"); err != nil {
		panic(err)
	}
	routerx.RouterConfig.OptDefalt()

	if *seedOnly {
		if err := seedAdmin(context.Background(), "admin", "admin123", "admin"); err != nil {
			logx.Logger.Errorf("seed failed: %v", err)
			os.Exit(1)
		}

		logx.Logger.Info("seed admin account created")
		return
	}

	// 登录时校验账号密码，account.LoadLoginAccount 会把密文解密后比对
	routerx.InitAuthMiddleware(account.LoadLoginAccount)

	jwtMW, err := routerx.AuthMiddleware()
	if err != nil {
		panic(err)
	}

	routerx.RegisterGroup(
		routerx.NewGroup("/api/v1").AddAPIGroup(
			routerx.NewGroup("/auth").AddAPI(
				&routerx.API{Method: "POST", Path: "/login", Handler: routerx.Login},
				&routerx.API{Method: "POST", Path: "/logout", Handler: routerx.Logout},
				&routerx.API{Method: "POST", Path: "/refresh", Handler: jwtMW.RefreshHandler},
			),
			routerx.NewGroup("/accounts").
				AddMiddlewares(jwtMW.MiddlewareFunc(), auditlog.Middleware.MiddlewareFunc()).
				AddAPI(
					&routerx.API{Method: "GET", Path: "", Handler: listAccounts},
					&routerx.API{Method: "POST", Path: "", Handler: createAccount},
					&routerx.API{Method: "GET", Path: "/me", Handler: routerx.GetAccountInfoHandler},
					&routerx.API{Method: "GET", Path: "/:id", Handler: getAccount},
					&routerx.API{Method: "PUT", Path: "/:id", Handler: updateAccount},
					&routerx.API{Method: "DELETE", Path: "/:id", Handler: deleteAccount},
				),
			routerx.NewGroup("/audit-logs").AddMiddlewares(jwtMW.MiddlewareFunc()).AddAPI(
				&routerx.API{Method: "GET", Path: "", Handler: auditlog.ListAuditLog},
				&routerx.API{Method: "GET", Path: "/export", Handler: auditlog.ExportAuditLog},
			),
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	utilsx.Notify(cancel)

	if err := routerx.StartRouter(ctx, routerx.NewRouter()); err != nil {
		panic(err)
	}
}
