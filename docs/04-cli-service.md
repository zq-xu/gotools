# 04 示例服务 C：命令行工具

同一套库也能支撑非 HTTP 的进程入口。本库自己的 CLI 就是一个例子：`cmd/gotools` 只有 5 行，命令实现全部放在 `internal/cli`。

## 本库 CLI 的结构

```
gotools/
├── cmd/gotools/main.go      # package main: cobra.CheckErr(cli.NewCommand().Execute())
└── internal/cli/
    ├── root.go              # 根命令 + 全局 --loglevel
    ├── format/              # gotools format
    └── excel/gantt/         # gotools excel gantt
```

`main.go` 全文：

```go
package main

import (
	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/internal/cli"
)

func main() {
	cobra.CheckErr(cli.NewCommand().Execute())
}
```

`internal/` 的语义是「只在本模块内可见」：外部项目无法 import 这些实现，所以库可以随时重构命令层而不构成破坏性变更。**写自己的 CLI 时不要试图 import 它**，照抄模式即可。

## 写一个自己的 CLI

```go
package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/logx"
)

var (
	cfgFile  string
	logLevel string
)

func main() {
	root := &cobra.Command{
		Use:   "demo",
		Short: "示例命令行工具",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return logx.SetLoggerLevel(logLevel)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("nothing to do, try: demo greet")
		},
	}

	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	root.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "log level")
	root.AddCommand(newGreetCommand())

	cobra.CheckErr(root.Execute())
}

func newGreetCommand() *cobra.Command {
	var times int

	cmd := &cobra.Command{
		Use:   "greet [name]",
		Short: "打招呼",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "world"
			if len(args) > 0 {
				name = args[0]
			}

			for i := 0; i < times; i++ {
				logx.Logger.Infof("hello, %s (%d/%d)", name, i+1, times)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&times, "times", "n", 1, "说几次")
	return cmd
}
```

约定与取舍：

- **`RunE` 优先于 `Run`**，错误一路抛到 `main` 由 `cobra.CheckErr` 统一退出，避免各处 `os.Exit`。
- **全局 flag 用 `PersistentFlags()`**，子命令 flag 用 `Flags()`；`cfgFile` / `logLevel` 这类每个命令都关心的放全局。
- **CLI 不一定要 `configx.Setup`**：纯本地命令（重排 import、生成文件）不需要数据库，强行 Setup 反而要求一份 config.yaml。需要库能力时再装配，并且注意「import 了 storex 就必须有 database 段」。
- 需要访问数据库的子命令，把 `configx.Setup(cfgFile)` 放在该子命令的 `PreRunE` 里，不要放在根命令。

## 内置命令参考

| 命令 | 作用 | 关键参数 |
| --- | --- | --- |
| `gotools --loglevel debug` | 设置全局日志级别 | 根命令持久 flag，默认 `info` |
| `gotools format` | 重排 Go 文件 `import (` 块 | `--dir` 目录（默认当前目录）、`--file` 单文件 |
| `gotools excel gantt plan.xlsx` | 在 Excel 里生成甘特图 | 位置参数：xlsx 路径 |

`format` 的分组规则：按 `go.mod` 的 module path，`import` 中含该前缀的归「本模块」组、含 `.` 的归第三方组、其余归标准库组，组间留空行，最后过一遍 `go/format`。它只重排 `import (` 块，不碰其它代码。

`excel gantt` 的输入约定：Sheet1 首行列名必须是 `Routes / GroupNo. / StartDate / EndDate`，日期形如 `2026/1/1`；结果写到同一文件的 Sheet2。

## 构建与安装

```bash
# 本库的 CLI
make build                                  # -> bin/gotools
go install github.com/zq-xu/gotools/cmd/gotools@latest

# 你的 CLI
go build -o bin/demo ./cmd/demo
```

因为入口在 `cmd/<name>`，`go install` 的路径也必须带上 `/cmd/<name>`；直接 `go install ...@latest` 会失败——模块根不再是 `package main`。

相关：[05 配置与日志](05-config-and-logging.md) 讲装配在非 HTTP 进程里怎么用。
