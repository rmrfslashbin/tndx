package ops

import (
	"github.com/rmrfslashbin/tndx/subcmds/ops/events"
	"github.com/spf13/cobra"
)

var (
	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Version: "v2021.12.28-00",
	}
)

func init() {
	RootCmd.AddCommand(
		events.RootCmd,
	)
}
