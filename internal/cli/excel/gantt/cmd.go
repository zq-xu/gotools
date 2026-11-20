package gantt

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

// Excel.Sheet1 example:
//
//	A          B            C          D
//
// Routes	GroupNo.	StartDate	EndDate
// Italian	GE29121623	2026/1/1	2026/1/3
func NewCommand() *cobra.Command {
	ganttCmd := &cobra.Command{
		Use:     "gantt",
		Short:   `generate gantt for excel`,
		Example: "gotools excel gantt aaa.xlsx",
		Long:    `operate excel`,
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return eris.Errorf("Empty args. Should set excel file path as the arg.")
			}

			filePath := args[0]
			return Gantt(filePath)
		},
	}

	return ganttCmd
}
