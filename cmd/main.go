package main

import (
	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/cmd/commands"
)

func main() {
	cobra.CheckErr(commands.NewCommand().Execute())
}
