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
