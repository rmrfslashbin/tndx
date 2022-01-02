package tndx

import (
	"github.com/rmrfslashbin/tndx/subcmds/tndx/tweets"
	"github.com/rmrfslashbin/tndx/subcmds/tndx/users"
	"github.com/spf13/cobra"
)

var (
	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Version: "v2022.01.02-00",
	}
)

func init() {
	RootCmd.AddCommand(
		tweets.RootCmd,
		users.RootCmd,
	)
}
