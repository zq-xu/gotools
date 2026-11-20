package main

import (
	"github.com/spf13/cobra"

	"github.com/zq-xu/gotools/internal/cli"
)

func main() {
	cobra.CheckErr(cli.NewCommand().Execute())
}
