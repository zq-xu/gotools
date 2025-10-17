package format

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/logx"
)

func NewFormatCommand() *cobra.Command {
	var dir, goFile string
	formatCmd := &cobra.Command{
		Use:     "format",
		Short:   `Format golang code`,
		Example: "gotools format",
		Long:    `Format golang code`,
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cobra.CheckErr(InitGoProject())

			// Validate flags
			if dir != "." && goFile != "" {
				return fmt.Errorf("cannot specify both --dir and --file")
			}

			// Decide which function to call
			if goFile != "" {
				logx.Logger.Infof("Start Go format file %s", goFile)
				defer logx.Logger.Info("Finish Go format")
				return FormatGoCodeInFile(goFile)
			}

			// default: use dir
			logx.Logger.Infof("Start Go format in dir %s", dir)
			defer logx.Logger.Info("Finish Go format")
			return FormatGoCodeInDir(dir)
		},
	}

	formatCmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "the golang code dir")
	formatCmd.PersistentFlags().StringVarP(&goFile, "file", "f", "", "the golang code file")
	return formatCmd
}
