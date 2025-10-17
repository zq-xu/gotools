package commands

import (
	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/cmd/commands/format"
	"github.com/zq-xu/gotools/logx"
)

var (
	logLevel string
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gotools",
		Short: `gotools`,
		Long:  `gotools`,
		Args:  cobra.ArbitraryArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(logx.SetLoggerLevel(logLevel))
			logx.Logger.Infof("LogLevel: %s", logLevel)
		},
		Run: func(cmd *cobra.Command, args []string) {
			logx.Logger.Info("Hello gotools")
		},
	}

	cmd.AddCommand(format.NewFormatCommand())
	cmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", logx.DefaultLogrusLogLevel.String(), "the log level")
	return cmd
}
