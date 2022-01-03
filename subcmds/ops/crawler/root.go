package crawler

import (
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/glue"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
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
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "", "dotenv path")

	cmdStatus.PersistentFlags().StringVarP(&flags.crawlerName, "crawler", "", "", "crawler name")

	cmdRun.PersistentFlags().StringVarP(&flags.crawlerName, "crawler", "", "", "crawler name")

	RootCmd.AddCommand(
		cmdList,
		cmdStatus,
		cmdRun,
	)
}

func setup() {
	if flags.dotenvPath == "" {
		// get platform specific user config directory
		configHome, err := os.UserConfigDir()
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("could not get user config directory and dotenv file not set")
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path.Join(configHome, "tndx"))
		viper.AddConfigPath(".")
	} else {
		flags.dotenvPath = path.Clean(flags.dotenvPath)
		viper.SetConfigFile(flags.dotenvPath)
		if _, err := os.Stat(flags.dotenvPath); err != nil {
			log.WithFields(logrus.Fields{
				"path":  flags.dotenvPath,
				"error": err,
			}).Fatal("unable to load dotenv")
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(logrus.Fields{
			"path": flags.dotenvPath,
			"err":  err,
		}).Fatal("failed to read dotenv file")
	}

	aws_region := viper.GetString("AwsRegion")

	if aws_region == "" {
		log.Fatal("AwsRegion not set in yaml config file")
	}

	if flags.crawlerName == "" {
		aws_profile := viper.GetString("AwsProfile")
		crawlerName := viper.GetString("CrawlerName")

		if aws_profile == "" {
			log.Fatal("AwsProfile not set in yaml config file")
		}
		if crawlerName == "" {
			log.Fatal("CrawlerName not set in yaml config file")
		}

		params := ssmparams.NewSSMParams(
			ssmparams.SetRegion(aws_region),
			ssmparams.SetProfile(aws_profile),
			ssmparams.SetLogger(log),
		)

		outputs, err := params.GetParams([]string{
			crawlerName,
		})

		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("failed to get params")
		}

		if len(outputs.InvalidParameters) > 0 {
			log.WithFields(logrus.Fields{
				"InvalidParameters": outputs.InvalidParameters,
			}).Fatal("invalid parameters")
		}

		flags.crawlerName = outputs.Params[crawlerName].(string)
	}

	crawler = glue.NewCrawler(
		glue.SetLogger(log),
		glue.SetRegion(aws_region),
		glue.SetCrawlerName(flags.crawlerName),
	)
}
