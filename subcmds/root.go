package subcmds

import (
	"fmt"
	"os"
	"path"

	"github.com/rmrfslashbin/ssmparams"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/rmrfslashbin/tndx/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type cliFlags struct {
	loglevel       string
	userid         int64
	screenname     string
	nextCursor     int64
	sinceid        int64
	maxid          int64
	backwards      bool
	storageDriver  string
	databaseDriver string
	localRootPath  string
	dotenvPath     string
	qEntities      bool
	runnerName     string
	friends        bool
	followers      bool
	favorites      bool
	timeline       bool
	all            bool
	none           bool
}

// service stores drivers and clients
type services struct {
	twitterClient *service.Config
	storage       storage.StorageDriver
	db            database.DatabaseDriver
	queue         *queue.Config
}

var (
	flags cliFlags
	svc   services
	log   *logrus.Logger

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Version: "v2021.12.04-00",
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

			setup(cmd)
		},
	}

	cmdLookup = &cobra.Command{
		Use:   "user",
		Short: "lookup user by userid or screenname",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.userid == 0 && flags.screenname == "" {
				cmd.Usage()
				log.Fatal("Missing userid or screenname")
				os.Exit(1)
			}

			if err := RunUserCmd(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdFriends = &cobra.Command{
		Use:   "friends",
		Short: "fetch the user's friends",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.userid == 0 && flags.screenname == "" {
				cmd.Usage()
				log.Fatal("Missing userid or screenname")
				os.Exit(1)
			}

			if err := RunFriendsCmd(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdFollowers = &cobra.Command{
		Use:   "followers",
		Short: "fetch the user's followers",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.userid == 0 && flags.screenname == "" {
				cmd.Usage()
				log.Fatal("Missing userid or screenname")
				os.Exit(1)
			}

			if err := RunFollowersCmd(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdFavorites = &cobra.Command{
		Use:   "favorites",
		Short: "fetch the user's favorites",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunFavoritesCmd(); err != nil {
				if flags.userid == 0 && flags.screenname == "" {
					cmd.Usage()
					log.Fatal("Missing userid or screenname")
					os.Exit(1)
				}

				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdTimeline = &cobra.Command{
		Use:   "timeline",
		Short: "fetch the user's timeline",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.userid == 0 && flags.screenname == "" {
				cmd.Usage()
				log.Fatal("Missing userid or screenname")
				os.Exit(1)
			}

			if err := RunTimelineCmd(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdDDB = &cobra.Command{
		Use:   "ddb",
		Short: "dynamodb operations",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flags.databaseDriver != "ddb" {
				cmd.Usage()
				log.Fatal("database driver must be 'ddb' for DDB operations")
				os.Exit(1)
			}
		},
	}

	cmdDDBRunner = &cobra.Command{
		Use:   "runner",
		Short: "dynamodb operations for runners",
	}

	cmdDDBRunnerSet = &cobra.Command{
		Use:   "set",
		Short: "configure a user to a runner",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flags.userid == 0 && flags.screenname == "" {
				cmd.Usage()
				log.Fatal("Missing userid or screenname")
				os.Exit(1)
			}

			if flags.followers || flags.friends || flags.favorites || flags.timeline {
				if flags.all {
					cmd.Usage()
					log.Fatal("all flag cannot be used with followers, friends, favorites or timeline")
					os.Exit(1)
				}
			}

			if flags.all && flags.none {
				cmd.Usage()
				log.Fatal("all and none flags cannot be used together")
				os.Exit(1)
			}

			setup(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunDDBRunnerSet(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdDDBRunnerList = &cobra.Command{
		Use:   "list",
		Short: "list runners",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setup(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunDDBRunnerList(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	cmdDDBRunnerDel = &cobra.Command{
		Use:   "del",
		Short: "delete a user from a runner",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flags.userid == 0 && flags.screenname == "" {
				cmd.Usage()
				log.Fatal("Missing userid or screenname")
				os.Exit(1)
			}
			setup(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunDDBRunnerDel(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}
)

// init sets up the CLI and flags
func init() {
	// Set the log level
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	RootCmd.PersistentFlags().StringVarP(&flags.loglevel, "loglevel", "", "info", "[error|warn|info|debug|trace]")
	RootCmd.PersistentFlags().StringVarP(&flags.storageDriver, "storage", "", "", "[local|s3]")
	RootCmd.PersistentFlags().StringVarP(&flags.databaseDriver, "database", "", "", "[sqlite|ddb]")
	RootCmd.PersistentFlags().StringVarP(&flags.localRootPath, "localrootpath", "", "./data", "local root path")
	RootCmd.PersistentFlags().StringVarP(&flags.dotenvPath, "dotenv", "", "./.env", "dotenv path")

	cmdFriends.PersistentFlags().Int64VarP(&flags.nextCursor, "nextcursor", "", 0, "next cursor (to overrride database return)")
	cmdFriends.PersistentFlags().BoolVarP(&flags.backwards, "backwards", "", false, "fetch backwards")
	cmdFriends.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "user id")
	cmdFriends.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screen name")

	cmdFollowers.PersistentFlags().Int64VarP(&flags.nextCursor, "nextcursor", "", 0, "next cursor (to overrride database return)")
	cmdFollowers.PersistentFlags().BoolVarP(&flags.backwards, "backwards", "", false, "fetch backwards")
	cmdFollowers.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "user id")
	cmdFollowers.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screen name")

	cmdFavorites.PersistentFlags().Int64VarP(&flags.maxid, "maxid", "", 0, "max id (to overrride database return)")
	cmdFavorites.PersistentFlags().Int64VarP(&flags.sinceid, "sinceid", "", 0, "since id (to overrride database return)")
	cmdFavorites.PersistentFlags().BoolVarP(&flags.backwards, "backwards", "", false, "fetch backwards")
	cmdFavorites.PersistentFlags().BoolVarP(&flags.qEntities, "queueentities", "", false, "send media entities to SQS")
	cmdFavorites.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "user id")
	cmdFavorites.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screen name")

	cmdTimeline.PersistentFlags().Int64VarP(&flags.maxid, "maxid", "", 0, "max id (to overrride database return)")
	cmdTimeline.PersistentFlags().Int64VarP(&flags.sinceid, "sinceid", "", 0, "since id (to overrride database return)")
	cmdTimeline.PersistentFlags().BoolVarP(&flags.backwards, "backwards", "", false, "fetch backwards")
	cmdTimeline.PersistentFlags().BoolVarP(&flags.qEntities, "queueentities", "", false, "send media entities to SQS")
	cmdTimeline.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "user id")
	cmdTimeline.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screen name")

	cmdDDBRunnerSet.PersistentFlags().StringVarP(&flags.runnerName, "runner", "", "", "runner name")
	cmdDDBRunnerSet.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "user id")
	cmdDDBRunnerSet.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screen name")
	cmdDDBRunnerSet.PersistentFlags().BoolVarP(&flags.favorites, "favorites", "", false, "favorites")
	cmdDDBRunnerSet.PersistentFlags().BoolVarP(&flags.followers, "followers", "", false, "followers")
	cmdDDBRunnerSet.PersistentFlags().BoolVarP(&flags.friends, "friends", "", false, "friends")
	cmdDDBRunnerSet.PersistentFlags().BoolVarP(&flags.timeline, "timeline", "", false, "timeline")
	cmdDDBRunnerSet.PersistentFlags().BoolVarP(&flags.all, "all", "", false, "set all (favorties, followers, friends, and timelien")
	cmdDDBRunnerSet.PersistentFlags().BoolVarP(&flags.none, "none", "", false, "unset/disable runners")
	cmdDDBRunnerSet.MarkPersistentFlagRequired("runner")

	cmdDDBRunnerList.PersistentFlags().StringVarP(&flags.runnerName, "runner", "", "", "runner name")
	cmdDDBRunnerList.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "user id")
	cmdDDBRunnerList.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screen name")
	cmdDDBRunnerList.MarkPersistentFlagRequired("runner")

	cmdDDBRunnerDel.PersistentFlags().StringVarP(&flags.runnerName, "runner", "", "", "runner name")
	cmdDDBRunnerDel.PersistentFlags().Int64VarP(&flags.userid, "userid", "", 0, "user id")
	cmdDDBRunnerDel.PersistentFlags().StringVarP(&flags.screenname, "screenname", "", "", "screen name")
	cmdDDBRunnerDel.MarkPersistentFlagRequired("runner")

	RootCmd.AddCommand(
		cmdLookup,
		cmdFriends,
		cmdFollowers,
		cmdFavorites,
		cmdTimeline,
		cmdDDB,
	)

	cmdDDB.AddCommand(
		cmdDDBRunner,
	)

	cmdDDBRunner.AddCommand(
		cmdDDBRunnerSet,
		cmdDDBRunnerList,
		cmdDDBRunnerDel,
	)
}

func setup(cmd *cobra.Command) {
	flags.dotenvPath = path.Clean(flags.dotenvPath)
	if _, err := os.Stat(flags.dotenvPath); err != nil {
		log.Fatal("Unable to find dotenv file")
		os.Exit(1)
	}

	validStorage := map[string]bool{
		"local": true,
		"s3":    true,
	}

	if _, ok := validStorage[flags.storageDriver]; !ok {
		cmd.Usage()
		log.Fatal("Invalid or missing storage driver. Should be 'local' or 's3'")
		os.Exit(1)
	}

	validDatabase := map[string]bool{
		"sqlite": true,
		"ddb":    true,
	}

	if _, ok := validDatabase[flags.databaseDriver]; !ok {
		cmd.Usage()
		log.Fatal("Invalid or missing database driver. Should be 'sqlite' or 'ddb'")
		os.Exit(1)
	}

	viper.SetConfigFile(flags.dotenvPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	twitterApiKey, _ := viper.Get("TWITTER_API_KEY_PARAM").(string)
	twitterApiSecret, _ := viper.Get("TWITTER_API_SECRET_PARAM").(string)
	s3Bucket, _ := viper.Get("TNDX_S3_BUCKET_PARAM").(string)
	s3Region, _ := viper.Get("TNDX_S3_REGION_PARAM").(string)
	ddbTable, _ := viper.Get("TNDX_DDB_TABLE_PARAM").(string)
	ddbRunnerTable, _ := viper.Get("TNDX_DDB_RUNNER_TABLE_PARAM").(string)
	ddbRegion, _ := viper.Get("TNDX_DDB_REGION_PARAM").(string)

	// awscli profile name
	awsProfile := "default"
	awsRegion := "us-east-1"

	// Set up a new ssmparams client
	params, err := ssmparams.New(
		ssmparams.SetProfile(awsProfile),
		ssmparams.SetRegion(awsRegion),
	)
	if err != nil {
		panic(err)
	}

	if flags.qEntities {
		sqsEntitiesURL, _ := viper.Get("TNDX_SQS_ENTITIES_URL_PARAM").(string)
		outputs, err := params.GetParams([]string{sqsEntitiesURL, s3Bucket})
		if err != nil {
			log.Fatal(err)
		}
		svc.queue = queue.NewSQS(
			queue.SetLogger(log),
			queue.SetSQSURL(outputs.Parameters[sqsEntitiesURL].(string)),
			queue.SetS3Bucket(outputs.Parameters[s3Bucket].(string)),
		)
	}

	if flags.databaseDriver == "sqlite" {
		// get platform specific user config directory
		confighome, _ := os.UserConfigDir()
		dbDir := path.Join(confighome, "tndx")
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatal(fmt.Sprintf("Unable to mkdir %s :: %s", dbDir, err.Error()))
			os.Exit(1)
		}
		configDB := path.Join(dbDir, "tndx.db")
		svc.db = database.NewSqliteDatabase(
			database.SetDatabaseFilename(configDB),
			database.SetSqliteLogger(log),
		)
	} else if flags.databaseDriver == "ddb" {
		outputs, err := params.GetParams([]string{ddbTable, ddbRegion, ddbRunnerTable})
		if err != nil {
			log.Fatal(err)
		}

		svc.db = database.NewDDB(
			database.SetDDBLogger(log),
			database.SetDDBTable(outputs.Parameters[ddbTable].(string)),
			database.SetDDBRunnerTable(outputs.Parameters[ddbRunnerTable].(string)),
			database.SetDDBRegion(outputs.Parameters[ddbRegion].(string)),
		)

	}

	outputs, err := params.GetParams([]string{twitterApiKey, twitterApiSecret, s3Bucket, s3Region})
	if err != nil {
		log.Fatal(err)
	}

	var storageDriver storage.StorageDriver
	if flags.storageDriver == "local" {
		storageDriver = storage.NewLocalStorage(storage.SetRootPath(flags.localRootPath))
	} else {

		storageDriver = storage.NewS3Storage(
			storage.SetS3Bucket(outputs.Parameters[s3Bucket].(string)),
			storage.SetS3Region(outputs.Parameters[s3Region].(string)),
		)

	}
	svc.storage = storageDriver

	svc.twitterClient = service.New(
		service.SetConsumerKey(outputs.Parameters[twitterApiKey].(string)),
		service.SetConsumerSecret(outputs.Parameters[twitterApiSecret].(string)),
		service.SetLogger(log),
	)
}
