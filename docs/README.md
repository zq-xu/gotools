# gotools 文档

按「要做什么」组织，不按包组织。每个示例都是完整文件，可以整块复制进新项目。

## 阅读路径

| 你的处境 | 建议顺序 |
| --- | --- |
| 第一次接触这个库 | [01 快速开始](01-quickstart.md) → [06 数据访问](06-database.md) → [07 HTTP 层](07-http.md) |
| 要写一个带登录的后端 | 01 → [02 REST 服务](02-rest-api-service.md) → [08 认证与审计](08-auth-and-audit.md) |
| 要写后台任务 / 消费队列 | 01 → [03 Worker 服务](03-worker-service.md) → [05 配置](05-config-and-logging.md) |
| 要写命令行工具 | [04 CLI](04-cli-service.md) → [05 配置](05-config-and-logging.md) |
| 只是查某个包的用法 | 直接跳 [05](05-config-and-logging.md) ~ [10](10-websocket.md) 的参考章节 |

## 文档清单

| 文件 | 类型 | 内容 |
| --- | --- | --- |
| [01-quickstart.md](01-quickstart.md) | 教程 | 20 行配置 + 40 行代码跑通一个服务 |
| [02-rest-api-service.md](02-rest-api-service.md) | 示例服务 A | 账号 CRUD + JWT + 审计，接近真实后端骨架 |
| [03-worker-service.md](03-worker-service.md) | 示例服务 B | 队列 + 固定 Worker 池的后台任务进程 |
| [04-cli-service.md](04-cli-service.md) | 示例服务 C | cobra 命令行工具，以及本库 CLI 的结构 |
| [05-config-and-logging.md](05-config-and-logging.md) | 参考 | 配置装配机制、注册表、日志 |
| [06-database.md](06-database.md) | 参考 | Model 约定、Database 接口、事务、分页与排序 |
| [07-http.md](07-http.md) | 参考 | 引擎装配、路由组 DSL、中间件链、优雅关闭 |
| [08-auth-and-audit.md](08-auth-and-audit.md) | 参考 | JWT 认证、角色校验、审计日志与导出 |
| [09-crypto.md](09-crypto.md) | 参考 | AES 加解密、密钥来源与轮换 |
| [10-websocket.md](10-websocket.md) | 参考 | Hub 模式连接管理与广播 |

## 三个示例服务的关系

```
02 REST 服务  ──┐
03 Worker     ──┼─→ 共用 01 的 config.yaml 结构与 05 的装配方式
04 CLI        ──┘
```

三者互不依赖，可以分别作为新项目的起点。02 覆盖读写的完整链路，03 覆盖长驻进程与并发任务，04 覆盖进程入口与命令分发。

## 术语约定

- **配置段**：`config.yaml` 顶层的一个 key，例如 `database:`；与一个注册项一一对应。
- **装配**：`configx.Setup` 遍历注册表，把每个配置段 unmarshal 到对应结构体并调用初始化函数的过程。
- **注册**：包在 `init()` 里调用 `configx.Register(name, &cfg, fn)`，因此 import 即生效。
- **handler**：只做「解析请求 → 调数据层 → 写响应」，不承载业务规则之外的逻辑。

示例默认 MySQL 5.7+ 可用、库名 `gotools_demo`、字符集 `utf8mb4`。代码以仓库 `main` 分支为准，示例中的 Go 代码均已随仓库编译验证。
