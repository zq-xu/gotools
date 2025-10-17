package main

import (
	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/cmd"
)

func main() {
	cobra.CheckErr(cmd.NewCommand().Execute())
}
