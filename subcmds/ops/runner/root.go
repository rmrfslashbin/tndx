package runner

import (
	"os"
	"path"

	"github.com/rmrfslashbin/ssmparams"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type Flags struct {
	loglevel   string
	userid     int64
	screenname string
	dotenvPath string
	runner     string
	friends    bool
	followers  bool
	favorites  bool
	timeline   bool
	user       bool
	all        bool
	none       bool
}

// service stores drivers and clients
type services struct {
	twitterClient *service.Config
	db            database.DatabaseDriver
}

var (
	flags Flags
	log   *logrus.Logger
	svc   services

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "runner",
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
		Short: "list runner entires",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			if flags.screenname == "" && flags.userid == 0 && !flags.all {
				cmd.Usage()
				log.Fatalf("--screenname or --userid must be set if not using --all")
			}

			if flags.all && flags.screenname != "" || flags.userid != 0 {
				cmd.Usage()
				log.Fatalf("--screenname or --userid must not be set if using --all")
			}

			if flags.runner == "" {
				cmd.Usage()
				log.Fatalf("runner must be set")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunRunnerList(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdSet = &cobra.Command{
		Use:   "set",
		Short: "set a runner user",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			if flags.screenname == "" && flags.userid == 0 {
				cmd.Usage()
				log.Fatalf("--screenname or --userid must be set")
			}

			if flags.all && flags.none {
				cmd.Usage()
				log.Fatalf("--all and --none cannot be set together")
			}

			if flags.all && (flags.friends || flags.followers || flags.favorites || flags.timeline || flags.user) {
				cmd.Usage()
				log.Fatalf("--all cannot be set with --friends, --followers, --favorites, --timeline, or --user")
			}

			if flags.none && (flags.friends || flags.followers || flags.favorites || flags.timeline || flags.user) {
				cmd.Usage()
				log.Fatalf("--none cannot be set with --friends, --followers, --favorites, --timeline, or --user")
			}

			if !flags.none && !flags.all && !flags.friends && !flags.followers && !flags.favorites && !flags.timeline && !flags.user {
				cmd.Usage()
				log.Fatalf("--friends, --followers, --favorites, --timeline, --user, --all or --none must be set")
			}

			if flags.runner == "" {
				cmd.Usage()
				log.Fatalf("runner must be set")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunRunnerSet(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdDel = &cobra.Command{
		Use:   "del",
		Short: "delete a runner user",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			if flags.screenname == "" && flags.userid == 0 {
				cmd.Usage()
				log.Fatalf("--screenname or --userid must be set")
			}

			if flags.runner == "" {
				cmd.Usage()
				log.Fatalf("runner must be set")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunRunnerDel(); err != nil {
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
	RootCmd.PersistentFlags().StringVarP(&flags.runner, "runner", "", "", "runner")

	cmdList.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screenname")
	cmdList.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "userid")
	cmdList.PersistentFlags().BoolVarP(&flags.all, "all", "", false, "all")

	cmdSet.PersistentFlags().StringVarP(&flags.runner, "runner", "", "", "runner")
	cmdSet.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screenname")
	cmdSet.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "userid")
	cmdSet.PersistentFlags().BoolVarP(&flags.friends, "friends", "", false, "friends")
	cmdSet.PersistentFlags().BoolVarP(&flags.followers, "followers", "", false, "followers")
	cmdSet.PersistentFlags().BoolVarP(&flags.favorites, "favorites", "", false, "favorites")
	cmdSet.PersistentFlags().BoolVarP(&flags.timeline, "timeline", "", false, "timeline")
	cmdSet.PersistentFlags().BoolVarP(&flags.user, "user", "", false, "user")
	cmdSet.PersistentFlags().BoolVarP(&flags.all, "all", "", false, "all")
	cmdSet.PersistentFlags().BoolVarP(&flags.none, "none", "", false, "none")

	cmdDel.PersistentFlags().StringVarP(&flags.runner, "runner", "", "", "runner")
	cmdDel.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screenname")
	cmdDel.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "userid")

	RootCmd.AddCommand(
		cmdList,
		cmdSet,
		cmdDel,
	)
}

func setup() {
	flags.dotenvPath = path.Clean(flags.dotenvPath)
	if _, err := os.Stat(flags.dotenvPath); err != nil {
		log.WithFields(logrus.Fields{
			"path":  flags.dotenvPath,
			"error": err,
		}).Fatal("unalbe to load dotenv")
	}

	viper.SetConfigFile(flags.dotenvPath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(logrus.Fields{
			"path": flags.dotenvPath,
			"err":  err,
		}).Fatal("failed to read dotenv file")
	}

	aws_region := viper.GetString("AWS_REGION")
	aws_profile := viper.GetString("AWS_PROFILE")
	ddb_table_prefix := viper.GetString("DDB_TABLE_PERFIX")
	twitter_api_key := viper.GetString("TWITTER_API_KEY")
	twitter_api_secret := viper.GetString("TWITTER_API_SECRET")

	if aws_region == "" {
		log.Fatal("AWS_REGION not set")
	}
	if aws_profile == "" {
		log.Fatal("AWS_PROFILE not set")
	}
	if ddb_table_prefix == "" {
		log.Fatal("DDB_TABLE_PERFIX not set")
	}
	if twitter_api_key == "" {
		log.Fatal("TWITTER_API_KEY not set")
	}
	if twitter_api_secret == "" {
		log.Fatal("TWITTER_API_SECRET not set")
	}

	// Set up a new ssmparams client
	params, err := ssmparams.New(
		ssmparams.SetProfile(aws_profile),
		ssmparams.SetRegion(aws_region),
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("error fetching ssm params")
	}

	outputs, err := params.GetParams([]string{
		ddb_table_prefix,
		twitter_api_key,
		twitter_api_secret,
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(outputs.InvalidParameters) > 0 {
		log.WithFields(logrus.Fields{
			"InvalidParameters": outputs.InvalidParameters,
		}).Fatal("invalid parameters")
	}

	svc.twitterClient = service.New(
		service.SetConsumerKey(outputs.Parameters[twitter_api_key].(string)),
		service.SetConsumerSecret(outputs.Parameters[twitter_api_secret].(string)),
		service.SetLogger(log),
	)

	svc.db = database.NewDDB(
		database.SetDDBLogger(log),
		database.SetDDBTablePrefix(outputs.Parameters[ddb_table_prefix].(string)),
		database.SetDDBRegion(aws_region),
	)
}
