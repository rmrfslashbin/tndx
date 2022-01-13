package exports

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/rmrfslashbin/tndx/pkg/database"
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
	tableArn   string
	s3Bucket   string
	s3Prefix   string
	format     string
	exportArn  string
}

type Services struct {
	twitter *service.Config
	db      *database.DDBDriver
}

var (
	flags *Flags
	log   *logrus.Logger
	svc   *Services

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "export",
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

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "start a ddb table export to s3",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch strings.ToUpper(flags.format) {
			case "DYNAMODB_JSON":
				return nil
			case "ION":
				return nil
			default:
				return errors.New("invalid --format expecting dynamodb_json or ion")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := runStartDDBExport(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "get the status of a ddb table export",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runStatusDDBExport(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	flags = &Flags{}
	svc = &Services{}
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	RootCmd.PersistentFlags().StringVarP(&flags.loglevel, "loglevel", "", "info", "[error|warn|info|debug|trace]")
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "", "dotenv path")

	startCmd.Flags().StringVarP(&flags.tableArn, "table-arn", "a", "", "arn of the table to export")
	startCmd.Flags().StringVarP(&flags.s3Bucket, "s3bucket", "b", "", "s3 bucket to export to")
	startCmd.Flags().StringVarP(&flags.s3Prefix, "s3prefix", "p", "", "s3 prefix to export to")
	startCmd.Flags().StringVarP(&flags.format, "format", "f", "json", "format to export to (DynamoDB_JSON, ION)")
	startCmd.MarkFlagRequired("table-arn")
	startCmd.MarkFlagRequired("s3-bucket")

	statusCmd.Flags().StringVarP(&flags.exportArn, "export-arn", "e", "", "arn of the export to get status for")
	statusCmd.MarkFlagRequired("export-arn")

	RootCmd.AddCommand(
		startCmd,
		statusCmd,
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
	params := ssmparams.NewSSMParams(
		ssmparams.SetRegion(aws_region),
		ssmparams.SetLogger(log),
	)

	outputs, err := params.GetParams([]string{
		ddb_table_prefix,
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

	svc.db = database.NewDDB(
		database.SetDDBLogger(log),
		database.SetDDBTablePrefix(outputs.Params[ddb_table_prefix].(string)),
		database.SetDDBRegion(aws_region),
	)

}
