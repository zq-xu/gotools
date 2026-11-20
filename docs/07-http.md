# 07 HTTP 层

`routerx` 负责引擎装配、路由注册与生命周期；`apix` 负责把「分页列表 / 详情 / 更新 / 删除」这类高度重复的 handler 写成泛型函数。

## 引擎装配

`routerx.NewRouter()` 返回的引擎自带以下中间件，顺序固定：

| 顺序 | 中间件 | 行为 |
| --- | --- | --- |
| 1 | `gin.Recovery()` | panic 转 500 |
| 2 | 请求日志 | 只打印**非 GET/OPTIONS**，并跳过 `/health` |
| 3 | `gzip` | 默认压缩级别 |
| 4 | CORS | 回显请求的 `Origin`，允许 credentials 与常用方法；`OPTIONS` 直接 204 |

另外两件事在装配时一次做完：

- `gin.SetMode(gin.ReleaseMode)`，控制台彩色输出关闭；
- `r.MaxMultipartMemory = 100 << 20`（100 MB）；
- 未匹配路由返回 `404 {"api":"not found"}`，与业务 404 区分得开。

请求日志只在写操作出现，是刻意的：GET 量大且无副作用，排障时看业务日志更有用。

## 路由组

不用 `r.GET(...)` 直接挂，而是先声明 `APIGroup` 再统一注册——好处是路径、中间件、子组都声明在一处，且能在 `NewRouter()` 之前完成收集。

```go
type APIGroup struct {
	RelativePath string
	Middlewares  []gin.HandlerFunc
	APIs         []*API
	Groups       []*APIGroup
}

type API struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}
```

```go
routerx.RegisterGroup(
	routerx.NewGroup("/api/v1").AddAPIGroup(
		routerx.NewGroup("/projects").
			AddMiddlewares(jwtMW.MiddlewareFunc()).
			AddAPI(
				&routerx.API{Method: "GET", Path: "", Handler: listProjects},
				&routerx.API{Method: "GET", Path: "/:id", Handler: getProject},
			),
	),
)
```

约定：

- `Path: ""` 表示直接落在组路径上（`/api/v1/projects`）；
- 子组路径拼接，中间件**继承但不重复**——父组中间件对子组生效；
- `:id` 统一用 `routerx.GetID(ctx)` 取，参数名固定 `id`；需要 int64 时 `routerx.GetIDInt64FromToken(ctx)` 取的是**token 里的用户 ID**，不是路径参数；
- `RegisterGroup` 是全局的，必须在 `NewRouter()` 之前调用。

## 健康检查

```go
GET /health  →  {}
```

想接入自己的探活逻辑（检查 DB、检查下游）就替换实现：

```go
routerx.SetHealthHandler(func(ctx *gin.Context) {
	if err := storex.GormDB(ctx).Exec("SELECT 1").Error; err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"db": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"db": "ok"})
})
```

## 错误响应

约定是 handler 里只写一行：

```go
if ei != nil {
	ctx.JSON(ei.StatusCode(), errorx.NewResponseBody(ei))
	return
}
```

`apix` / `getx` 的内置 handler 已全部走 `errorx.NewResponseBody`，响应体统一为 `{"error": "<message>"}`：

```bash
curl -s 'localhost:8080/api/v1/accounts?pageNum=abc'
# {"error":"PageNum is invalid. strconv.Atoi: parsing \"abc\": invalid syntax"}
```

自己的 handler 里同样传 `errorx.NewResponseBody(ei)`；只有裸 `error` 时用 `errorx.ResponseBody{Error: err.Error()}`，前后拿到的结构就一致了。

## 生命周期

```go
ctx, cancel := context.WithCancel(context.Background())
utilsx.Notify(cancel)   // SIGINT / SIGTERM → cancel

err := routerx.StartRouter(ctx, routerx.NewRouter())
```

- 监听地址来自 `router.RouterConfig.Host` / `Port`，所以 **`OptDefalt()` 要在 `StartRouter` 前调用**；
- `DisableTLS=false` 时用 `certPath` / `keyPath` 起 TLS，开发环境务必设 `disableTLS: true`；
- `ctx` 取消后走 `Shutdown`，最多等 **3 秒**，超时强制关闭；
- 返回的错误只包含真正的启动失败（`http.ErrServerClosed` 被过滤掉），所以 `if err != nil { panic(err) }` 是安全写法。

## 泛型 handler（apix）

`apix` 把四类固定形态的接口做成泛型函数，业务只需要给模型和回调：

| handler | 语义 | 成功状态码 | 失败 |
| --- | --- | --- | --- |
| `ListHandler[T]` | 分页列表 | 200 `PageResponse` | 400 |
| `ListByRawGormHandler[T]` | 同上，但查询用 `*gorm.DB` 自己拼 | 200 | 400 |
| `UpdateHandler[T, P]` | 事务内取对象 → 回调改 → 保存 | **201** | 404 / 400 |
| `UpdateByRawGormHandler[T, P]` | 同上，但读取与保存走 raw gorm | 201 | 404 / 400 |
| `DeleteHandler[T]` | 事务内取对象 → 校验 → 删除 | **204** | 400 |
| `getx.GetHandler[T, R]` | 取单个对象并转响应 | 200 | 404 |

设计上的取舍：

- **没有 CreateHandler**：创建时的查重、初始化、默认值每个模型都不同，强行泛化只会把业务规则塞进回调里；
- **Delete 幂等**：对象不存在也返回 204，调用方不需要区分「删掉了」和「本来就没有」；
- **Update 返回值固定 201**，与 HTTP 语义不完全一致，但库内保持统一，前端按状态码判断时注意；
- `DeleteHandler` 的 `validateFn` 允许为 nil（历史版本传 nil 会 panic）。

`ListByRawGormHandler` 的总数与列表用同一套条件计算：`queryFn` 附加的 `WHERE`、查询条件、模糊搜索都计入，只排除分页：

```go
apix.ListByRawGormHandler[Order](ctx,
	func(g *gorm.DB, p *typesx.ListParams) *gorm.DB {
		return g.Where("status = ?", 1)
	},
	apix.DefaultTransListObjToResp[Order, OrderResponse],
)
```

此时 `count` 只统计 `status=1` 的行，与列表一致。

## 完整示例

路由表怎么和这些 handler 组合，见 [02 REST 服务](02-rest-api-service.md)；认证与审计中间件见 [08](08-auth-and-audit.md)。
