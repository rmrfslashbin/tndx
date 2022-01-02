package crawler

import (
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/glue"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type Flags struct {
	loglevel    string
	dotenvPath  string
	crawlerName string
}

var (
	flags   Flags
	log     *logrus.Logger
	crawler *glue.Config

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "crawler",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set the log level
			switch flags.loglevel {
			case "error":
				log.SetLevel(logrus.ErrorLevel)
			case "warn":
				log.SetLevel(logrus.WarnLevel)
			case "info":
				log.SetLevel(logrus.InfoLevel)
			case "debug":
				log.SetLevel(logrus.DebugLevel)
			case "trace":
				log.SetLevel(logrus.TraceLevel)
			default:
				log.SetLevel(logrus.InfoLevel)
			}
			setup()
		},
	}

	cmdList = &cobra.Command{
		Use:   "list",
		Short: "list crawlers",
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runCrawlerList(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdStatus = &cobra.Command{
		Use:   "status",
		Short: "show crawler status",
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runCrawlerStatus(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdRun = &cobra.Command{
		Use:   "run",
		Short: "run crwaler",
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runRunCrawler(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	flags = Flags{}
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	RootCmd.PersistentFlags().StringVarP(&flags.loglevel, "loglevel", "", "info", "[error|warn|info|debug|trace]")
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "./.env", "dotenv path")

	cmdStatus.PersistentFlags().StringVarP(&flags.crawlerName, "crawler", "", "", "crawler name")
	cmdStatus.MarkFlagRequired("crawler")

	cmdRun.PersistentFlags().StringVarP(&flags.crawlerName, "crawler", "", "", "crawler name")
	cmdRun.MarkFlagRequired("crawler")

	RootCmd.AddCommand(
		cmdList,
		cmdStatus,
		cmdRun,
	)
}

func setup() {
	flags.dotenvPath = path.Clean(flags.dotenvPath)
	if _, err := os.Stat(flags.dotenvPath); err != nil {
		log.WithFields(logrus.Fields{
			"path":  flags.dotenvPath,
			"error": err,
		}).Fatal("unable to load dotenv")
	}

	viper.SetConfigFile(flags.dotenvPath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(logrus.Fields{
			"path": flags.dotenvPath,
			"err":  err,
		}).Fatal("failed to read dotenv file")
	}

	aws_region := viper.GetString("AWS_REGION")

	if aws_region == "" {
		log.Fatal("AWS_REGION not set")
	}

	crawler = glue.NewCrawler(
		glue.SetLogger(log),
		glue.SetRegion(aws_region),
		glue.SetCrawlerName(flags.crawlerName),
	)
}
