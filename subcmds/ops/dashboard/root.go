package dashboard

import (
	"fmt"
	"os"
	"path"
	"text/tabwriter"

	"github.com/rmrfslashbin/tndx/pkg/events"
	"github.com/rmrfslashbin/tndx/pkg/glue"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags struct contains settings for the root command
type Flags struct {
	loglevel   string
	dotenvPath string
}

var (
	flags Flags
	log   *logrus.Logger
	e     *events.Config
	q     *queue.Config
	c     *glue.Config

	// rootCmd is the Viper root command
	RootCmd = &cobra.Command{
		Use: "dashboard",
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
		Run: func(cmd *cobra.Command, args []string) {
			run()
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

	e = events.NewEvents(
		events.SetLogger(log),
		events.SetRegion(aws_region),
		events.SetProfile(aws_profile),
	)

	q = queue.NewSQS(
		queue.SetLogger(log),
		queue.SetRegion(aws_region),
		queue.SetProfile(aws_profile),
		queue.SetSQSURL(outputs.Params[sqs_queue_url].(string)),
	)

	c = glue.NewCrawler(
		glue.SetLogger(log),
		glue.SetRegion(aws_region),
		glue.SetProfile(aws_profile),
		glue.SetCrawlerName("tndx-rmrfslashbin-tweets"),
	)
}

func run() {
	if ret, err := q.GetAttribs(); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to get queue attributes")
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Queue Attrib\tValue")
		fmt.Fprintf(w, "Queue ARN\t%s\n", ret["QueueArn"])
		fmt.Fprintf(w, "ApproximateNumberOfMessages\t%s\n", ret["ApproximateNumberOfMessages"])
		fmt.Fprintf(w, "ApproximateNumberOfMessagesNotVisible\t%s\n", ret["ApproximateNumberOfMessagesNotVisible"])
		fmt.Fprintf(w, "ApproximateNumberOfMessagesDelayed\t%s\n", ret["ApproximateNumberOfMessagesDelayed"])

		/*
			for k, v := range ret {
				fmt.Fprintf(w, "%s\t%s\n", k, v)
			}
		*/
		w.Flush()
		fmt.Println()
	}

	if rules, err := e.List(); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to list events")
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Event\tDescription\tRate\tStatus")
		for _, rule := range rules.Rules {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", *rule.Name, *rule.Description, *rule.ScheduleExpression, rule.State)
		}
		w.Flush()
		fmt.Println()
	}

	if ret, err := c.GetCrawlerData(); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to get crawler data")
	} else {
		fmt.Println("Crawler Status")
		fmt.Printf("Name:                   %s\n", *ret.Crawler.Name)
		fmt.Printf("State:                  %s\n", *&ret.Crawler.State)
		fmt.Printf("Elapsed Time:           %d\n", ret.Crawler.CrawlElapsedTime)
		fmt.Printf("Last Crawl Status:      %v\n", ret.Crawler.LastCrawl.Status)
		fmt.Printf("Last Crawl Error:       %v\n", ret.Crawler.LastCrawl.ErrorMessage)
		fmt.Printf("Last Crawl Start time:  %v\n", ret.Crawler.LastCrawl.StartTime)
		fmt.Println()
	}
}
