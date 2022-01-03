package events

import (
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/events"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type Flags struct {
	loglevel   string
	dotenvPath string
	ruleName   string
}

var (
	flags Flags
	log   *logrus.Logger
	evnts *events.Config

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "events",
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
		Short: "list event rules",
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runEventsList(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdDisable = &cobra.Command{
		Use:   "disable",
		Short: "disable event rule",
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runEventDisable(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdEnable = &cobra.Command{
		Use:   "enable",
		Short: "enable event rule",
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runEventEnable(); err != nil {
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

	cmdDisable.Flags().StringVarP(&flags.ruleName, "rule", "", "", "rule name")
	cmdDisable.MarkFlagRequired("rule")

	cmdEnable.Flags().StringVarP(&flags.ruleName, "rule", "", "", "rule name")
	cmdEnable.MarkFlagRequired("rule")

	RootCmd.AddCommand(
		cmdList,
		cmdDisable,
		cmdEnable,
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
	aws_profile := viper.GetString("AwsProfile")
	ddb_table_prefix := viper.GetString("DDBTablePrefix")
	twitter_api_key := viper.GetString("TwitterApiKey")
	twitter_api_secret := viper.GetString("TwitterApiSecret")
	sqs_queue_url := viper.GetString("SQSQueueUrl")

	if aws_region == "" {
		log.Fatal("AwsRegion not set in yaml config file")
	}
	if aws_profile == "" {
		log.Fatal("AwsProfile not set in yaml config file")
	}
	if ddb_table_prefix == "" {
		log.Fatal("DDBTablePrefix not set in yaml config file")
	}
	if twitter_api_key == "" {
		log.Fatal("TwitterApiKey not set in yaml config file")
	}
	if twitter_api_secret == "" {
		log.Fatal("TwitterApiSecret not set in yaml config file")
	}
	if sqs_queue_url == "" {
		log.Fatal("SQSQueueUrl not set in yaml config file")
	}

	params := ssmparams.NewSSMParams(
		ssmparams.SetRegion(aws_region),
		ssmparams.SetProfile(aws_profile),
		ssmparams.SetLogger(log),
	)

	outputs, err := params.GetParams([]string{
		ddb_table_prefix,
		twitter_api_key,
		twitter_api_secret,
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

	evnts = events.NewEvents(
		events.SetLogger(log),
		events.SetRegion(aws_region),
	)
}
