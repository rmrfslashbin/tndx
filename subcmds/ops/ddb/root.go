package ddb

import (
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/rmrfslashbin/tndx/subcmds/ops/ddb/params"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type Flags struct {
	loglevel   string
	dotenvPath string
	userid     int64
	tweetid    int64
	friendid   int64
	followid   int64
	screenname string
	json       bool
	yaml       bool
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
		Use: "ddb",
	}

	favoriteCmd = &cobra.Command{
		Use:   "favorites",
		Short: "fetch favorites for a user",
		PreRun: func(cmd *cobra.Command, args []string) {
			if flags.tweetid != 0 && (flags.userid != 0 || flags.screenname != "") {
				log.Fatal("--tweetid and --userid/--screenname are mutually exclusive")
			}
			if flags.tweetid == 0 && flags.userid == 0 && flags.screenname == "" {
				log.Fatal("must specify --tweetid or at least one --userid/--screenname")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runDDBFavorites(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	friendsCmd = &cobra.Command{
		Use:   "friends",
		Short: "fetch friends for a user",
		PreRun: func(cmd *cobra.Command, args []string) {
			if flags.friendid != 0 && (flags.userid != 0 || flags.screenname != "") {
				log.Fatal("--friendid and --userid/--screenname are mutually exclusive")
			}
			if flags.friendid == 0 && flags.userid == 0 && flags.screenname == "" {
				log.Fatal("must specify --friendid or at least one --userid/--screenname")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runDDBFriends(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	followersCmd = &cobra.Command{
		Use:   "followers",
		Short: "fetch followers for a user",
		PreRun: func(cmd *cobra.Command, args []string) {
			if flags.followid != 0 && (flags.userid != 0 || flags.screenname != "") {
				log.Fatal("--followid and --userid/--screenname are mutually exclusive")
			}
			if flags.followid == 0 && flags.userid == 0 && flags.screenname == "" {
				log.Fatal("must specify --followid or at least one --userid/--screenname")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			setup()
			if err := runDDBFollowers(); err != nil {
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

	RootCmd.PersistentFlags().StringVarP(&flags.loglevel, "loglevel", "l", "info", "log level (debug, info, warn, error, fatal, panic)")

	favoriteCmd.Flags().StringVarP(&flags.screenname, "screenname", "s", "", "screenname of user to fetch favorites for")
	favoriteCmd.Flags().Int64VarP(&flags.userid, "userid", "u", 0, "userid of user to fetch favorites for")
	favoriteCmd.Flags().Int64VarP(&flags.tweetid, "tweetid", "t", 0, "tweetid to fetch favorites for")
	favoriteCmd.Flags().BoolVarP(&flags.json, "json", "j", false, "output in json format")
	favoriteCmd.Flags().BoolVarP(&flags.yaml, "yaml", "y", false, "output in yaml format")

	followersCmd.Flags().StringVarP(&flags.screenname, "screenname", "s", "", "screenname of user to fetch followers for")
	followersCmd.Flags().Int64VarP(&flags.userid, "userid", "u", 0, "userid of user to fetch followers for")
	followersCmd.Flags().Int64VarP(&flags.followid, "followid", "f", 0, "followid to fetch followers for")
	followersCmd.Flags().BoolVarP(&flags.json, "json", "j", false, "output in json format")
	followersCmd.Flags().BoolVarP(&flags.yaml, "yaml", "y", false, "output in yaml format")

	friendsCmd.Flags().StringVarP(&flags.screenname, "screenname", "s", "", "screenname of user to fetch friends for")
	friendsCmd.Flags().Int64VarP(&flags.userid, "userid", "u", 0, "userid of user to fetch friends for")
	friendsCmd.Flags().Int64VarP(&flags.friendid, "friendid", "f", 0, "friendid to fetch frineds  for")
	friendsCmd.Flags().BoolVarP(&flags.json, "json", "j", false, "output in json format")
	friendsCmd.Flags().BoolVarP(&flags.yaml, "yaml", "y", false, "output in yaml format")

	RootCmd.AddCommand(
		params.RootCmd,
		favoriteCmd,
		followersCmd,
		friendsCmd,
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
