package tndx

import (
	"github.com/rmrfslashbin/tndx/subcmds/tndx/tweets"
	"github.com/rmrfslashbin/tndx/subcmds/tndx/users"
	"github.com/spf13/cobra"
)

var (
	Version string
	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Version: Version,
	}
)

func init() {
	RootCmd.AddCommand(
		tweets.RootCmd,
		users.RootCmd,
	)
}
