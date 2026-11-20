# 08 认证与审计

`routerx/auth` 把 [gin-jwt](https://github.com/appleboy/gin-jwt) 包成本库约定；`routerx/auditlog` 负责给写操作留痕并支持导出。

## 认证

### 装配顺序

```go
// 1. 告诉中间件「怎么校验账号密码」
routerx.InitAuthMiddleware(account.LoadLoginAccount)

// 2. 拿到中间件实例（RefreshHandler 等需要它）
jwtMW, err := routerx.AuthMiddleware()
if err != nil {
	panic(err)
}

// 3. 注册登录路由与受保护路由组
```

`InitAuthMiddleware` 的入参是 `func(ctx, username, password) (auth.AuthAccount, errorx.ErrorInfo)`。任何实现了下面五个 getter 的模型都能用：

```go
type AuthAccount interface {
	GetID() string
	GetUsername() string
	GetRoles() string
	GetName() string
	GetStatus() string
}
```

`account.Account` 已经实现，所以内置账号直接用 `account.LoadLoginAccount`；接第三方用户系统就自己写一个 loader。

### 端点

| 路由 | 方法 | 行为 |
| --- | --- | --- |
| `/auth/login` | POST | 校验 `{username, password}`，成功返回 **201** `{token, expire}` |
| `/auth/logout` | POST | 返回 201（无 body）；纯客户端行为，token 不会失效 |
| `/auth/refresh` | POST | 挂在 `jwtMW.RefreshHandler` 上，用旧 token 换新 token |
| `/accounts/me` | GET | `routerx.GetAccountInfoHandler`，回显当前用户信息 |

登录失败时 gin-jwt 走 `Unauthorized`：普通错误与 token 过期都返回 `{"errorCode":401,"errorMessage":"..."}`，形状一致。

### token 里有什么

```go
AuthAccountIDToken       = "auth_account_id"
AuthAccountNameToken     = "auth_account_name"
AuthAccountUsernameToken = "auth_account_username"
AuthUserRolesToken       = "auth_account_roles"    // 逗号分隔的字符串
AuthAccountStatusToken   = "auth_account_status"
```

取值：

```go
info := routerx.GetUserInfoFromToken(ctx)
// &AccountInfoResponse{ID, Username, Name, Roles []string, Status}
```

- **只能在已过 JWT 中间件的路由上调用**。未认证时 claims 为空，函数返回全空结构体（不会 panic），但拿它做判断一定会得到空值。
- `Roles` 由逗号分隔字符串拆出来，`"admin,operator"` → `["admin", "operator"]`。
- 需要 int64 的用户 ID：`routerx.GetIDInt64FromToken(ctx)`。

### 角色校验

```go
cfg := getx.NewDefaultConfig[account.Account, account.Account]()
cfg.ValidateFn = getx.ValidateAuthInRoles("admin", "operator")
getx.GetHandler[account.Account, account.Account](ctx, cfg)
```

命中任意一个角色放行，都不命中返回 **403**；没有 claims 时返回 **401**。其它 handler（`apix.ListHandler` 等）没有 ValidateFn 参数，需要自己在 handler 开头判断：

```go
if roles := routerx.GetUserInfoFromToken(ctx).Roles; !slices.Contains(roles, "admin") {
	ctx.JSON(http.StatusForbidden, "permission denied")
	return
}
```

### JWT 参数通过 Option 配置

```go
routerx.InitAuthMiddleware(account.LoadLoginAccount,
	routerx.WithKey(os.Getenv("JWT_SECRET")), // 生产必须设置，缺省会用开发用默认密钥并打 warn
	routerx.WithTimeout(time.Hour),
	routerx.WithMaxRefresh(time.Hour*100),
)

Key:           o.Key        // 默认 "secret key"（不安全，仅限开发）
Timeout:       o.Timeout    // 默认 1 小时
MaxRefresh:    o.MaxRefresh // 默认 100 小时
TokenLookup:   "header: Authorization, query: token, cookie: jwt"
TokenHeadName: "Bearer"
```

不传 Option 时保持历史默认值并打印 warn 日志；`TokenLookup` / `TokenHeadName` 仍是固定值，需要别的取 token 方式就自己建 `jwt.GinJWTMiddleware`。

## 审计

### 挂载

```go
routerx.NewGroup("/projects").
	AddMiddlewares(jwtMW.MiddlewareFunc(), auditlog.Middleware.MiddlewareFunc()).
	AddAPI(...)
```

`auditlog.Middleware` 是包级实例（`*auditLogMiddleware`），通过 `MiddlewareFunc()` 变成 `gin.HandlerFunc`。它自己会把 `ModelAuditLog` 注册进自动迁移，不需要额外 `RegisterTable`。

### 记录规则

| 条件 | 结果 |
| --- | --- |
| 方法为 `GET` / `OPTIONS` | 跳过 |
| 方法为 `POST` / `PUT` / `DELETE` | 记录 |
| 响应状态码 ≥ 400 | 跳过（只留成功的写操作） |
| 落库失败 | 打 error 日志，不影响业务响应 |

落库用 `context.Background()`，请求被取消也不会丢审计。记录包含方法、路径、ClientIP、请求体、响应体、状态码。

`ModelAuditLog.UserID / UserName / UserAlias` 三个字段目前固定为空——库内取值逻辑被注释掉了。要关联操作人，就在自己的中间件里补写，或者改库。

### 敏感字段脱敏

默认记录原始请求体。涉及密码、手机号时按「方法 + 路径」注册自定义取体函数，它返回的字节会作为 `Body` 存档：

```go
auditlog.Middleware.RegisterBodyGenerateFn("POST", "/api/v1/accounts", func(ctx *gin.Context) []byte {
	body, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var m map[string]any
	_ = json.Unmarshal(body, &m)
	delete(m, "password")

	safe, _ := json.Marshal(m)
	return safe
})
```

### 查询与导出

两个 handler 是普通函数，自己挂路由：

```go
routerx.NewGroup("/audit-logs").AddMiddlewares(jwtMW.MiddlewareFunc()).AddAPI(
	&routerx.API{Method: "GET", Path: "", Handler: auditlog.ListAuditLog},
	&routerx.API{Method: "GET", Path: "/export", Handler: auditlog.ExportAuditLog},
)
```

- 列表：同样接受 `pageNum` / `pageSize` / `sortBy` / 精确过滤，返回 `PageResponse`，元素是 `ResponseOfAuditLog{id, createdAt, updatedAt, client_ip, url, method, body, message, status_code}`。
- 导出：`GET /export` 直接下载 `auditlog.xlsx`，列顺序取 `ExcelRowData` 的字段名（`UserID, UserName, ClientIP, Url, Method, Body, Message, StatusCode, CreatedAt`）。其中 `UserID` / `UserName` 列因上面说的原因恒为空。
- 两者都用同一个 `ListParams`，所以「导出哪一批」由 query 决定，没有单独的筛选 UI。

参考实现见 [02 REST 服务](02-rest-api-service.md)。
