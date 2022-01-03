package tweets

import (
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/queue"
	q "github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type Flags struct {
	loglevel   string
	dotenvPath string
	runner     string
	tweetids   []string
}

// service stores drivers and clients
type services struct {
	queue *q.Config
}

var (
	flags     *Flags
	log       *logrus.Logger
	svc       services
	bootstrap *queue.Bootstrap

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "tweets",
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

	cmdProcess = &cobra.Command{
		Use:   "process",
		Short: "process one or more tweets",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Parent().PersistentPreRun(cmd.Parent(), args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunTweetsProcess(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	flags = &Flags{}
	log = logrus.New()
	bootstrap = &queue.Bootstrap{}
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	RootCmd.PersistentFlags().StringVarP(&flags.loglevel, "loglevel", "", "info", "[error|warn|info|debug|trace]")
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "", "dotenv path")
	RootCmd.PersistentFlags().StringVarP(&flags.runner, "runner", "", "", "runner")

	cmdProcess.Flags().StringSliceVarP(&flags.tweetids, "tweetid", "", []string{}, "tweetid")
	cmdProcess.MarkFlagRequired("tweetid")

	RootCmd.AddCommand(
		cmdProcess,
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
	s3_bucket := viper.GetString("S3Bucket")
	sqs_queue_url := viper.GetString("SQSQueueUrl")
	tweet_delivery_stream := viper.GetString("TweetDeliveryStream")
	twitter_api_key := viper.GetString("TwitterApiKey")
	twitter_api_secret := viper.GetString("TwitterApiSecret")

	if aws_region == "" {
		log.Fatal("AWS_REGION not set")
	}
	if ddb_table_prefix == "" {
		log.Fatal("DDB_TABLE_PERFIX is required")
	}
	if s3_bucket == "" {
		log.Fatal("S3_BUCKET is required")
	}
	if sqs_queue_url == "" {
		log.Fatal("SQS_QUEUE_URL not set")
	}
	if tweet_delivery_stream == "" {
		log.Fatal("TWEET_DELIVERY_STREAM not set")
	}
	if twitter_api_key == "" {
		log.Fatal("TWITTER_API_KEY not set")
	}
	if twitter_api_secret == "" {
		log.Fatal("TWITTER_API_SECRET not set")
	}

	// Set up a new ssmparams client
	params := ssmparams.NewSSMParams(
		ssmparams.SetRegion(aws_region),
		ssmparams.SetProfile(aws_profile),
		ssmparams.SetLogger(log),
	)

	outputs, err := params.GetParams([]string{
		sqs_queue_url,
	})
	if err != nil {
		log.Fatal(err)
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

	bootstrap = &queue.Bootstrap{
		DDBTablePrefix:   ddb_table_prefix,
		DeliveryStream:   tweet_delivery_stream,
		SQSRunnerURL:     sqs_queue_url,
		S3Bucket:         s3_bucket,
		TwitterAPIKey:    twitter_api_key,
		TwitterAPISecret: twitter_api_secret,
	}
}
