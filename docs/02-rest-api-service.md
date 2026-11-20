# 02 示例服务 A：账号管理 REST API

一个接近真实后端的骨架：JWT 登录、账号 CRUD、分页查询、写操作审计。整个服务只拆成两个文件。

| 能力 | 由谁提供 |
| --- | --- |
| 启动装配 | `configx` |
| 路由表与中间件 | `routerx` |
| 登录与 token | `routerx` + `routerx/auth` |
| 账号模型与密码校验 | `storex/account` |
| 列表 / 详情 / 更新 / 删除 | `apix`、`apix/getx` |
| 写操作留痕与导出 | `routerx/auditlog` |

## 项目结构

```
demo/
├── config.yaml
├── main.go        # 装配 + 路由表
└── handlers.go    # 业务 handler
```

## 配置

```yaml
logs:
  logLevel: info

database:
  address: 127.0.0.1
  port: 3306
  username: root
  password: root
  database: gotools_demo
  logLevel: 4

router:
  host: 0.0.0.0
  port: 8080
  disableTLS: true

crypto:
  aesKey: "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4"   # base64 的 24 字节密钥
```

账号密码在库里以密文保存，`crypto` 段由 `storex/account` 注册并用于 `LoadLoginAccount` 解密；缺这一段登录一定失败。

## main.go：装配与路由表

```go
package main

import (
	"context"

	"github.com/zq-xu/gotools/configx"
	"github.com/zq-xu/gotools/routerx"
	"github.com/zq-xu/gotools/routerx/auditlog"
	"github.com/zq-xu/gotools/storex/account"
	"github.com/zq-xu/gotools/utilsx"
)

func main() {
	if err := configx.Setup("config.yaml"); err != nil {
		panic(err)
	}
	routerx.RouterConfig.OptDefalt()

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
```

几点值得注意：

- 组可以嵌套，`NewGroup("/api/v1").AddAPIGroup(NewGroup("/accounts")...)` 拼出 `/api/v1/accounts`；`Path: ""` 表示直接落在组路径上。
- 中间件按注册顺序执行，`jwtMW.MiddlewareFunc()` 必须在前把 claims 写进 ctx，后面的 `auditlog.Middleware.MiddlewareFunc()` 才能识别「谁做的操作」。
- `routerx.RegisterGroup` 是全局注册，必须在 `NewRouter()` 之前调用；`NewRouter()` 时统一挂到引擎。
- `/:id` 的取值统一走 `routerx.GetID(ctx)`，参数名固定 `id`。

## handlers.go

```go
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

// seedAdmin 是一次性初始化：给空库写入第一个可登录账号。
// 真实项目里通常做成子命令或迁移脚本，避免每次启动都执行。
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
```

每个 handler 的职责边界：

| handler | 做的事 | 不做什么 |
| --- | --- | --- |
| `createAccount` | 绑定入参、加密密码、查重、落库 | 不做角色校验（在路由层或 ValidateFn 里） |
| `listAccounts` | 交给 `apix.ListHandler`：解析分页/排序/过滤 → 查询计数 → 转响应 | 不拼 SQL |
| `getAccount` | 按 `id` 取详情并拷贝成响应 | 不判权限 |
| `updateAccount` | 事务内取原对象、应用 diff、保存 | 不做字段级校验（在 `ShouldBindJSON` 之外自己加） |
| `deleteAccount` | 事务内取对象、校验、删除 | 不做幂等之外的清理 |

`apix` 只提供 List / Get / Update / Delete 四个通用 handler；**Create 一律手写**，因为它的校验规则（查重、初始化、初始状态）每个模型都不同。

## 初始化第一个账号

第一次启动时 `accounts` 表是空的，没有账号可登录。在 `main` 里临时调一次：

```go
if err := seedAdmin(context.Background(), "admin", "admin123", "admin"); err != nil {
	panic(err)
}
```

跑通后删掉这段，或者把它改成 `seed` 子命令（见 [04 CLI](04-cli-service.md)）。

## 请求示例

登录（成功返回 201）：

```bash
curl -s -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}'
# {"token":"eyJhbGciOi...","expire":"2026-11-20T18:00:00+01:00"}
```

```bash
TOKEN='eyJhbGciOi...'
```

列表（分页、排序、精确过滤、模糊搜索）：

```bash
# 第 2 页，每页 10 条，按创建时间倒序
curl -s 'localhost:8080/api/v1/accounts?pageNum=2&pageSize=10&sortBy=created_at,desc' \
  -H "Authorization: Bearer $TOKEN"

# 按 username 精确过滤（key 用模型字段的 snake_case）
curl -s 'localhost:8080/api/v1/accounts?username=admin' -H "Authorization: Bearer $TOKEN"
```

详情、当前登录人、更新、删除：

```bash
curl -s localhost:8080/api/v1/accounts/123456789 -H "Authorization: Bearer $TOKEN"
curl -s localhost:8080/api/v1/accounts/me -H "Authorization: Bearer $TOKEN"

curl -s -X PUT localhost:8080/api/v1/accounts/123456789 \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Alan","roles":"admin,operator"}'

curl -s -o /dev/null -w '%{http_code}\n' -X DELETE localhost:8080/api/v1/accounts/123456789 \
  -H "Authorization: Bearer $TOKEN"
```

审计日志（只记录成功的写操作）：

```bash
curl -s 'localhost:8080/api/v1/audit-logs?pageNum=1&pageSize=20' -H "Authorization: Bearer $TOKEN"
curl -s -o auditlog.xlsx 'localhost:8080/api/v1/audit-logs/export' -H "Authorization: Bearer $TOKEN"
```

## 行为约定

| 端点 | 成功 | 失败 |
| --- | --- | --- |
| `POST /auth/login` | 201 `{token, expire}` | 401 账号或密码错误 |
| `GET /accounts` | 200 `PageResponse` | 400 分页参数非法 |
| `POST /accounts` | 201 账号对象 | 400 参数错误 / 409 用户名已存在 |
| `GET /accounts/:id` | 200 账号对象 | 404 不存在 |
| `PUT /accounts/:id` | **201**（沿用库内约定） | 404 不存在 / 400 参数或保存失败 |
| `DELETE /accounts/:id` | 204 无 body | 不存在时同样 204（幂等） |

`PUT` 返回 201 是 `apix.UpdateHandler` 的既有实现，前端按状态码做特判时注意。

## 加角色限制

`getx` 提供了角色校验，挂到详情类 handler 上：

```go
func getAccountAsAdmin(ctx *gin.Context) {
	cfg := getx.NewDefaultConfig[account.Account, account.Account]()
	cfg.ValidateFn = getx.ValidateAuthInRoles("admin", "operator")

	getx.GetHandler[account.Account, account.Account](ctx, cfg)
}
```

角色取自 token 的 `roles`（逗号分隔），命中任意一个即通过，否则 403。它只在**已过 JWT 中间件**的路由上有效——claims 是中间件写进 ctx 的。

## 已知限制

- `routerx/auth` 的 JWT `Key` 与 `Timeout` 目前写死，用于真实环境前需要改库或替换中间件。
- 审计记录里 `UserID / UserName` 字段目前固定为空（库内取值逻辑被注释掉了），要关联操作人需要自行补。
- `apix` 不提供 Create handler，也不提供批量操作。

下一步：同一套装配方式换成长驻后台任务，见 [03 Worker 服务](03-worker-service.md)。
