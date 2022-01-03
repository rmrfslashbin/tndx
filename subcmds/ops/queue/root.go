package queue

import (
	"os"
	"path"

	q "github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Services struct {
	queue *q.Config
}

// Flags struct contains settings for the root command
type Flags struct {
	loglevel   string
	dotenvPath string
	confirm    bool
}

var (
	flags Flags
	log   *logrus.Logger
	svc   *Services

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "queue",
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

	cmdPurge = &cobra.Command{
		Use:   "purge",
		Short: "purge queue",
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runQueuePurge(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	flags = Flags{}
	svc = &Services{}
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	RootCmd.PersistentFlags().StringVarP(&flags.loglevel, "loglevel", "", "info", "[error|warn|info|debug|trace]")
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "", "dotenv path")

	cmdPurge.Flags().BoolVarP(&flags.confirm, "confirm", "", false, "confirm purge")
	cmdPurge.MarkFlagRequired("confirm")
	RootCmd.AddCommand(cmdPurge)
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
	sqs_queue_url := viper.GetString("SQSQueueUrl")

	if aws_region == "" {
		log.Fatal("AwsRegion not set in yaml config file")
	}
	if aws_profile == "" {
		log.Fatal("AwsProfile not set in yaml config file")
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
		sqs_queue_url,
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

	svc.queue = q.NewSQS(
		q.SetLogger(log),
		q.SetRegion(aws_region),
		q.SetSQSURL(outputs.Params[sqs_queue_url].(string)),
	)
}
