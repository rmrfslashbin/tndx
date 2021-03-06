package database

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	rekognitionTypes "github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/sirupsen/logrus"
)

type Bits uint8

type DDBOption func(config *DDBDriver)

type DDBDriver struct {
	log                         *logrus.Logger
	driverName                  string
	tablePrefix                 string
	region                      string
	profile                     string
	favoritesTable              string
	favoritesTableGSITweetid    string
	friendsTable                string
	friendsTableGSIFriendid     string
	followersTable              string
	followersTableGSIFollowerid string
	runnerTable                 string
	mediaTable                  string
	paramsTable                 string
	db                          *dynamodb.Client
}
type TweetConfigQuery struct {
	UserID     int64
	SinceID    int64
	MaxID      int64
	LastUpdate int64
}

type CursoredTweetConfigQuery struct {
	UserID         int64
	PreviousCursor int64
	NextCursor     int64
}
type TweetsItem struct {
	Domain              string    `json:"Domain" yaml:"Domain"`
	UserID              int64     `json:"UserID" yaml:"UserID"`
	MaxID               int64     `json:"MaxID" yaml:"MaxID"`
	SinceID             int64     `json:"SinceID" yaml:"SinceID"`
	LastUpdate          int64     `json:"LastUpdate" yaml:"LastUpdate"`
	LastUpdateTimestamp time.Time `json:"LastUpdateTimestamp" yaml:"LastUpdateTimestamp"`
}

type UserToTweetLink struct {
	UserID  int64 `json:"UserID"`
	TweetID int64 `json:"TweetID"`
}

type UserToFollowerLink struct {
	UserID     int64 `json:"UserID"`
	FollowerID int64 `json:"FollowerID"`
}

type UserToFriendLink struct {
	UserID   int64 `json:"UserID"`
	FriendID int64 `json:"FriendID"`
}

type FavoritesItem struct {
	Domain              string    `json:"Domain" yaml:"Domain"`
	UserID              int64     `json:"UserID" yaml:"UserID"`
	MaxID               int64     `json:"MaxID" yaml:"MaxID"`
	SinceID             int64     `json:"SinceID" yaml:"SinceID"`
	LastUpdate          int64     `json:"LastUpdate" yaml:"LastUpdate"`
	LastUpdateTimestamp time.Time `json:"LastUpdateTimestamp" yaml:"LastUpdateTimestamp"`
}

type FollowersItem struct {
	Domain              string    `json:"Domain" yaml:"Domain"`
	UserID              int64     `json:"UserID" yaml:"UserID"`
	NextCursor          int64     `json:"NextCursor" yaml:"NextCursor"`
	PreviousCursor      int64     `json:"PreviousCursor" yaml:"PreviousCursor"`
	LastUpdate          int64     `json:"LastUpdate" yaml:"LastUpdate"`
	LastUpdateTimestamp time.Time `json:"LastUpdateTimestamp" yaml:"LastUpdateTimestamp"`
}

type FriendsItem struct {
	Domain              string    `json:"Domain" yaml:"Domain"`
	UserID              int64     `json:"UserID" yaml:"UserID"`
	NextCursor          int64     `json:"NextCursor" yaml:"NextCursor"`
	PreviousCursor      int64     `json:"PreviousCursor" yaml:"PreviousCursor"`
	LastUpdate          int64     `json:"LastUpdate" yaml:"LastUpdate"`
	LastUpdateTimestamp time.Time `json:"LastUpdateTimestamp" yaml:"LastUpdateTimestamp"`
}

type RunnerItem struct {
	RunnerName string `json:"RunnerName"`
	UserID     int64  `json:"UserID"`
	Flags      Bits   `json:"Flags"`
	LastUpdate int64  `json:"LastUpdate"`
}

type MediaItem struct {
	Bucket          string                             `json:"Bucket"`
	S3Key           string                             `json:"S3Key"`
	UserID          int64                              `json:"UserID"`
	TweetID         int64                              `json:"TweetID"`
	Faces           []rekognitionTypes.FaceDetail      `json:"Faces"`
	Labels          []rekognitionTypes.Label           `json:"Labels"`
	Moderation      []rekognitionTypes.ModerationLabel `json:"Moderation"`
	Text            []rekognitionTypes.TextDetection   `json:"Text"`
	FacesCount      int                                `json:"FacesCount"`
	LabelsCount     int                                `json:"LabelsCount"`
	ModerationCount int                                `json:"ModerationCount"`
	TextCount       int                                `json:"TextCount"`
}

const (
	F_favorites Bits = 1 << iota
	F_followers
	F_friends
	F_timeline
	F_user
)

func Set(b, flag Bits) Bits    { return b | flag }
func Clear(b, flag Bits) Bits  { return b &^ flag }
func Toggle(b, flag Bits) Bits { return b ^ flag }
func Has(b, flag Bits) bool    { return b&flag != 0 }

func NewDDB(opts ...func(*DDBDriver)) *DDBDriver {
	cfg := &DDBDriver{}
	cfg.driverName = "ddb"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.region == "" {
		cfg.region = os.Getenv("AWS_REGION")
	}

	c, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = cfg.region

		if cfg.profile != "" {
			o.SharedConfigProfile = cfg.profile
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(c)
	cfg.db = svc

	return cfg
}

func SetDDBProfile(profile string) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.profile = profile
	}
}

func SetDDBRegion(region string) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.region = region
	}
}

func SetDDBTablePrefix(tablePrefix string) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.tablePrefix = tablePrefix
		config.favoritesTable = tablePrefix + "favorites"
		config.favoritesTableGSITweetid = tablePrefix + "favorites-gsi-tweetid"
		config.friendsTable = tablePrefix + "friends"
		config.friendsTableGSIFriendid = tablePrefix + "friends-gsi-friendid"
		config.followersTable = tablePrefix + "followers"
		config.followersTableGSIFollowerid = tablePrefix + "followers-gsi-followerid"
		config.runnerTable = tablePrefix + "runners"
		config.mediaTable = tablePrefix + "media"
		config.paramsTable = tablePrefix + "parameters"
	}
}

func SetDDBLogger(logger *logrus.Logger) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.log = logger
	}
}

func (config *DDBDriver) DeleteMedia(mediaItem *MediaItem) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(config.mediaTable),
		Key: map[string]types.AttributeValue{
			"TweetID": &types.AttributeValueMemberN{Value: strconv.FormatInt(mediaItem.TweetID, 10)},
			"S3Key":   &types.AttributeValueMemberS{Value: mediaItem.S3Key},
		},
	}
	_, err := config.db.DeleteItem(context.TODO(), input)
	return err
}

func (config *DDBDriver) DeleteRunnerUser(params *RunnerItem) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(config.runnerTable),
		Key: map[string]types.AttributeValue{
			"RunnerName": &types.AttributeValueMemberS{Value: params.RunnerName},
			"UserID":     &types.AttributeValueMemberN{Value: strconv.FormatInt(params.UserID, 10)},
		},
	}
	_, err := config.db.DeleteItem(context.TODO(), input)
	return err
}

func (config *DDBDriver) GetDriverName() string {
	return config.driverName
}

func (config *DDBDriver) GetFavoritesByTweetId(tweetID int64) ([]*UserToTweetLink, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(config.favoritesTable),
		IndexName:              aws.String(config.favoritesTableGSITweetid),
		KeyConditionExpression: aws.String("TweetID = :TweetID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":TweetID": &types.AttributeValueMemberN{Value: strconv.FormatInt(tweetID, 10)},
		},
	}
	result, err := config.db.Query(context.TODO(), input)
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"error": err,
			"input": input,
		}).Error("Error querying favorite/users")
		return nil, err
	}

	results := []*UserToTweetLink{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetFavoritesByUserId(userID int64) ([]*UserToTweetLink, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(config.favoritesTable),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
		},
	}
	result, err := config.db.Query(context.TODO(), input)
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"error": err,
			"input": input,
		}).Error("Error querying user/favorites")
		return nil, err
	}

	results := []*UserToTweetLink{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetFollowersByFollowId(followID int64) ([]*UserToFollowerLink, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(config.followersTable),
		IndexName:              aws.String(config.followersTableGSIFollowerid),
		KeyConditionExpression: aws.String("FollowerID = :FollowerID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":FollowerID": &types.AttributeValueMemberN{Value: strconv.FormatInt(followID, 10)},
		},
	}
	result, err := config.db.Query(context.TODO(), input)
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"error": err,
			"input": input,
		}).Error("Error querying friend/users")
		return nil, err
	}

	results := []*UserToFollowerLink{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetFollowersByUserId(userID int64) ([]*UserToFollowerLink, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(config.followersTable),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
		},
	}
	result, err := config.db.Query(context.TODO(), input)
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"error": err,
			"input": input,
		}).Error("Error querying user/followers")
		return nil, err
	}

	results := []*UserToFollowerLink{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetFriendsByFriendId(friendID int64) ([]*UserToFriendLink, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(config.friendsTable),
		IndexName:              aws.String(config.friendsTableGSIFriendid),
		KeyConditionExpression: aws.String("FriendID = :FriendID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":FriendID": &types.AttributeValueMemberN{Value: strconv.FormatInt(friendID, 10)},
		},
	}
	result, err := config.db.Query(context.TODO(), input)
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"error": err,
			"input": input,
		}).Error("Error querying friend/users")
		return nil, err
	}

	results := []*UserToFriendLink{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetFriendsByUserId(userID int64) ([]*UserToFriendLink, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(config.friendsTable),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
		},
	}
	result, err := config.db.Query(context.TODO(), input)
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"error": err,
			"input": input,
		}).Error("Error querying user/friends")
		return nil, err
	}

	results := []*UserToFriendLink{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetFavoritesConfig(userID int64) (*FavoritesItem, error) {
	result, err := config.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(config.paramsTable),
		Key: map[string]types.AttributeValue{
			"UserID": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"Domain": &types.AttributeValueMemberS{Value: "favorites"},
		},
	})

	if err != nil {
		return nil, err
	}

	item := &FavoritesItem{}

	if result.Item == nil {
		return item, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetFollowersConfig(userID int64) (*FollowersItem, error) {
	result, err := config.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(config.paramsTable),
		Key: map[string]types.AttributeValue{
			"UserID": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"Domain": &types.AttributeValueMemberS{Value: "followers"},
		},
	})

	if err != nil {
		return nil, err
	}
	item := &FollowersItem{}

	if result.Item == nil {
		return item, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetFriendsConfig(userID int64) (*FriendsItem, error) {
	result, err := config.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(config.paramsTable),
		Key: map[string]types.AttributeValue{
			"UserID": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"Domain": &types.AttributeValueMemberS{Value: "friends"},
		},
	})

	if err != nil {
		return nil, err
	}

	item := &FriendsItem{}

	if result.Item == nil {
		return item, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetRunnerUsers(runnerUsers *RunnerItem) ([]*RunnerItem, error) {
	var input *dynamodb.QueryInput

	if runnerUsers.UserID == 0 {
		input = &dynamodb.QueryInput{
			TableName:              aws.String(config.runnerTable),
			KeyConditionExpression: aws.String("RunnerName = :RunnerName"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":RunnerName": &types.AttributeValueMemberS{Value: runnerUsers.RunnerName},
			},
		}
	} else {
		input = &dynamodb.QueryInput{
			TableName:              aws.String(config.runnerTable),
			KeyConditionExpression: aws.String("RunnerName = :RunnerName and UserID = :UserID"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":RunnerName": &types.AttributeValueMemberS{Value: runnerUsers.RunnerName},
				":UserID":     &types.AttributeValueMemberN{Value: strconv.FormatInt(runnerUsers.UserID, 10)},
			},
		}
	}

	result, err := config.db.Query(context.TODO(), input)
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"error": err,
			"input": input,
		}).Error("Error querying runner users")
		return nil, err
	}

	results := []*RunnerItem{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetTimelineConfig(userID int64) (*TweetsItem, error) {
	result, err := config.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(config.paramsTable),
		Key: map[string]types.AttributeValue{
			"UserID": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"Domain": &types.AttributeValueMemberS{Value: "tweets"},
		},
	})

	if err != nil {
		return nil, err
	}

	item := &TweetsItem{}

	if result.Item == nil {
		return item, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) PutFavorites(links []*UserToTweetLink) error {
	for _, link := range links {
		kvp, err := attributevalue.MarshalMap(link)
		if err != nil {
			return err
		}

		if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(config.favoritesTable),
			Item:      kvp,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (config *DDBDriver) PutFollowers(links []*UserToFollowerLink) error {
	for _, link := range links {
		kvp, err := attributevalue.MarshalMap(link)
		if err != nil {
			return err
		}

		if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(config.followersTable),
			Item:      kvp,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (config *DDBDriver) PutFriends(links []*UserToFriendLink) error {
	for _, link := range links {
		kvp, err := attributevalue.MarshalMap(link)
		if err != nil {
			return err
		}

		if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(config.friendsTable),
			Item:      kvp,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (config *DDBDriver) PutFavoritesConfig(query *TweetConfigQuery) error {
	now := time.Now()

	kvp, err := attributevalue.MarshalMap(&FavoritesItem{
		Domain:     "favorites",
		UserID:     query.UserID,
		MaxID:      query.MaxID,
		SinceID:    query.SinceID,
		LastUpdate: now.UnixMilli(),
	})
	if err != nil {
		return err
	}

	if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(config.paramsTable),
		Item:      kvp,
	}); err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutFollowersConfig(query *CursoredTweetConfigQuery) error {
	now := time.Now()
	kvp, err := attributevalue.MarshalMap(&FollowersItem{
		Domain:         "followers",
		UserID:         query.UserID,
		PreviousCursor: query.PreviousCursor,
		NextCursor:     query.NextCursor,
		LastUpdate:     now.UnixMilli(),
	})
	if err != nil {
		return err
	}
	if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.paramsTable),
	}); err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutFriendsConfig(query *CursoredTweetConfigQuery) error {
	now := time.Now()
	kvp, err := attributevalue.MarshalMap(&FriendsItem{
		Domain:         "friends",
		UserID:         query.UserID,
		PreviousCursor: query.PreviousCursor,
		NextCursor:     query.NextCursor,
		LastUpdate:     now.UnixMilli(),
	})
	if err != nil {
		return err
	}

	if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.paramsTable),
	}); err != nil {
		return err
	}
	return nil
}
func (config *DDBDriver) PutMedia(mediaItem *MediaItem) error {
	kvp, err := attributevalue.MarshalMap(mediaItem)
	if err != nil {
		return err
	}

	if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.mediaTable),
	}); err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutTimelineConfig(query *TweetConfigQuery) error {
	now := time.Now()
	kvp, err := attributevalue.MarshalMap(&TweetsItem{
		Domain:     "tweets",
		UserID:     query.UserID,
		MaxID:      query.MaxID,
		SinceID:    query.SinceID,
		LastUpdate: now.UnixMilli(),
	})
	if err != nil {
		return err
	}

	if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.paramsTable),
	}); err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutRunnerFlags(params *RunnerItem) error {
	now := time.Now()
	kvp, err := attributevalue.MarshalMap(&RunnerItem{
		RunnerName: params.RunnerName,
		UserID:     params.UserID,
		Flags:      params.Flags,
		LastUpdate: now.UnixMilli(),
	})
	if err != nil {
		return err
	}

	if _, err := config.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.runnerTable),
	}); err != nil {
		return err
	}
	return nil
}

type TableExportRequest struct {
	ClientToken string
	//DYNAMODB_JSON or ION
	ExportFormat string
	ExportTime   time.Time

	S3Bucket string
	S3Prefix string
	TableArn string

	exportFormat types.ExportFormat
}

func (config *DDBDriver) ExportTable(params *TableExportRequest) (*dynamodb.ExportTableToPointInTimeOutput, error) {
	switch strings.ToUpper(params.ExportFormat) {
	case "DYNAMODB_JSON":
		params.exportFormat = types.ExportFormatDynamodbJson
	case "ION":
		params.exportFormat = types.ExportFormatIon
	default:
		config.log.Info("defaulting to dynamodb json export format")
		params.exportFormat = types.ExportFormatDynamodbJson
	}

	return config.db.ExportTableToPointInTime(
		context.TODO(),
		&dynamodb.ExportTableToPointInTimeInput{
			S3Bucket:       aws.String(params.S3Bucket),
			TableArn:       aws.String(params.TableArn),
			ExportFormat:   params.exportFormat,
			ExportTime:     aws.Time(params.ExportTime),
			S3Prefix:       aws.String(params.S3Prefix),
			S3SseAlgorithm: types.S3SseAlgorithmAes256,
		},
	)
}

func (config *DDBDriver) ExportStatus(exportArn string) (*dynamodb.DescribeExportOutput, error) {
	return config.db.DescribeExport(
		context.TODO(),
		&dynamodb.DescribeExportInput{
			ExportArn: aws.String(exportArn),
		},
	)
}

func (config *DDBDriver) ExportList(tableArn string) (*dynamodb.ListExportsOutput, error) {
	return config.db.ListExports(
		context.TODO(),
		&dynamodb.ListExportsInput{
			TableArn: aws.String(tableArn),
		},
	)
}
