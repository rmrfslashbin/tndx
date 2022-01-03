package timeline

import (
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/kinesis"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type Flags struct {
	loglevel   string
	dotenvPath string
	userid     int64
	screenname string
	sinceid    int64
	maxid      int64
	count      int
}

type Services struct {
	twitter *service.Config
	kinesis *kinesis.Config
	queue   *queue.Config
}

var (
	flags     *Flags
	log       *logrus.Logger
	svc       *Services
	bootstrap *queue.Bootstrap

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "timeline",
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

	cmdIngest = &cobra.Command{
		Use:   "ingest",
		Short: "ping one or more tweets",
		PreRun: func(cmd *cobra.Command, args []string) {
			if flags.count > 200 {
				log.Warn("count is greater than 200, setting to 200")
				flags.count = 200
			}
			if flags.count < 1 {
				log.Warn("count is less than 1, setting to 1")
				flags.count = 1
			}
			if flags.userid == 0 && flags.screenname == "" {
				log.Fatal("userid or screenname is required")
			}
			if flags.userid != 0 && flags.screenname != "" {
				log.Fatal("userid and screenname are mutually exclusive")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runTimelineIngest(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	flags = &Flags{}
	svc = &Services{}
	bootstrap = &queue.Bootstrap{}

	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	RootCmd.PersistentFlags().StringVarP(&flags.loglevel, "loglevel", "", "info", "[error|warn|info|debug|trace]")
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "", "dotenv path")

	cmdIngest.Flags().Int64VarP(&flags.userid, "userid", "u", 0, "user id")
	cmdIngest.Flags().Int64VarP(&flags.sinceid, "sinceid", "s", 0, "since id")
	cmdIngest.Flags().Int64VarP(&flags.maxid, "maxid", "m", 0, "max id")
	cmdIngest.Flags().IntVarP(&flags.count, "count", "c", 200, "count")
	cmdIngest.Flags().StringVarP(&flags.screenname, "screenname", "n", "", "screen name")
	cmdIngest.MarkFlagRequired("sinceid")

	RootCmd.AddCommand(
		cmdIngest,
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

	params := ssmparams.NewSSMParams(
		ssmparams.SetRegion(aws_region),
		ssmparams.SetLogger(log),
	)

	outputs, err := params.GetParams([]string{
		ddb_table_prefix,
		s3_bucket,
		sqs_queue_url,
		tweet_delivery_stream,
		twitter_api_key,
		twitter_api_secret,
	})

	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "getParams",
			"error":  err.Error(),
		}).Fatal("error getting parameters.")
	}

	if len(outputs.InvalidParameters) > 0 {
		log.WithFields(logrus.Fields{
			"invalid_parameters": outputs.InvalidParameters,
		}).Fatal("invalid parameters")
	}

	svc.twitter = service.New(
		service.SetConsumerKey(outputs.Params[twitter_api_key].(string)),
		service.SetConsumerSecret(outputs.Params[twitter_api_secret].(string)),
		service.SetLogger(log),
	)

	svc.kinesis = kinesis.NewFirehose(
		kinesis.SetRegion(aws_region),
		kinesis.SetLogger(log),
		kinesis.SetProfile(aws_profile),
		kinesis.SetDeliveryStream(outputs.Params[tweet_delivery_stream].(string)),
	)

	svc.queue = queue.NewSQS(
		queue.SetLogger(log),
		queue.SetRegion(aws_region),
		queue.SetProfile(aws_profile),
		queue.SetSQSURL(outputs.Params[sqs_queue_url].(string)),
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
