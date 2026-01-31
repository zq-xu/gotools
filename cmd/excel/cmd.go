package excel

import (
	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/cmd/excel/gantt"
)

func NewExcelCommand() *cobra.Command {
	excelCmd := &cobra.Command{
		Use:     "excel",
		Short:   `operate excel`,
		Example: "gotools excel",
		Long:    `operate excel`,
		Args:    cobra.ArbitraryArgs,
	}

	excelCmd.AddCommand(gantt.NewCommand())
	return excelCmd
}
