package ddb

import (
	"github.com/rmrfslashbin/tndx/subcmds/ops/ddb/params"
	"github.com/spf13/cobra"
)

var (
	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "ddb",
	}
)

func init() {
	RootCmd.AddCommand(
		params.RootCmd,
	)
}
