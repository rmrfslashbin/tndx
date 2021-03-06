package ops

import (
	"github.com/rmrfslashbin/tndx/subcmds/ops/crawler"
	"github.com/rmrfslashbin/tndx/subcmds/ops/dashboard"
	"github.com/rmrfslashbin/tndx/subcmds/ops/ddb"
	"github.com/rmrfslashbin/tndx/subcmds/ops/events"
	"github.com/rmrfslashbin/tndx/subcmds/ops/queue"
	"github.com/rmrfslashbin/tndx/subcmds/ops/runner"
	"github.com/rmrfslashbin/tndx/subcmds/ops/tweets"
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
		events.RootCmd,
		dashboard.RootCmd,
		crawler.RootCmd,
		runner.RootCmd,
		queue.RootCmd,
		tweets.RootCmd,
		ddb.RootCmd,
	)
}
