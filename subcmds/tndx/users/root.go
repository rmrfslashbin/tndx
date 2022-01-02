package users

import (
	"os"
	"path"

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
	userids    []int64
	screenname []string
	json       bool
	yaml       bool
	basic      bool
}

type Services struct {
	twitter *service.Config
}

var (
	flags *Flags
	log   *logrus.Logger
	svc   *Services

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "users",
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

	cmdGet = &cobra.Command{
		Use:   "get",
		Short: "get one or more users",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(flags.userids) == 0 && len(flags.screenname) == 0 {
				log.Fatal("must specify at least one userid or screenname")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runUsersGet(); err != nil {
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
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "./.env", "dotenv path")

	cmdGet.Flags().Int64SliceVarP(&flags.userids, "userid", "u", []int64{}, "user id")
	cmdGet.Flags().StringSliceVarP(&flags.screenname, "screenname", "s", []string{}, "screen name")
	cmdGet.Flags().BoolVarP(&flags.json, "json", "j", false, "output in json format")
	cmdGet.Flags().BoolVarP(&flags.yaml, "yaml", "y", false, "output in yaml format")
	cmdGet.Flags().BoolVarP(&flags.basic, "basic", "b", false, "output in basic format")

	RootCmd.AddCommand(
		cmdGet,
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
	twitter_api_key := viper.GetString("TWITTER_API_KEY")
	twitter_api_secret := viper.GetString("TWITTER_API_SECRET")

	if aws_region == "" {
		log.Fatal("AWS_REGION not set")
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

}

func DedupInt64Slice(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func DedupStringSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
