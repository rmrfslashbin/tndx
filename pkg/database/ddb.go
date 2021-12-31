package database

import (
	"context"
	"os"
	"strconv"
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
	log            *logrus.Logger
	driverName     string
	tablePrefix    string
	region         string
	favoritesTable string
	friendsTable   string
	followersTable string
	runnerTable    string
	mediaTable     string
	paramsTable    string
	db             *dynamodb.Client
}
type TweetConfigQuery struct {
	UserID  int64
	SinceID int64
	MaxID   int64
}

type CursoredTweetConfigQuery struct {
	UserID         int64
	PreviousCursor int64
	NextCursor     int64
}
type TweetsItem struct {
	Domain     string `json:"Domain"`
	UserID     int64  `json:"UserID"`
	MaxID      int64  `json:"MaxID"`
	SinceID    int64  `json:"SinceID"`
	LastUpdate int64  `json:"LastUpdate"`
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
	Domain     string `json:"Domain"`
	UserID     int64  `json:"UserID"`
	MaxID      int64  `json:"MaxID"`
	SinceID    int64  `json:"SinceID"`
	LastUpdate int64  `json:"LastUpdate"`
}

type FollowersItem struct {
	Domain         string `json:"Domain"`
	UserID         int64  `json:"UserID"`
	NextCursor     int64  `json:"NextCursor"`
	PreviousCursor int64  `json:"PreviousCursor"`
	LastUpdate     int64  `json:"LastUpdate"`
}

type FriendsItem struct {
	Domain         string `json:"Domain"`
	UserID         int64  `json:"UserID"`
	NextCursor     int64  `json:"NextCursor"`
	PreviousCursor int64  `json:"PreviousCursor"`
	LastUpdate     int64  `json:"LastUpdate"`
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
		return nil
	})
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(c)
	cfg.db = svc

	return cfg
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
		config.friendsTable = tablePrefix + "friends"
		config.followersTable = tablePrefix + "followers"
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

func (config *DDBDriver) GetFavoritesConfig(userID int64) (*TweetConfigQuery, error) {
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

	item := &TweetConfigQuery{}

	if result.Item == nil {
		return item, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetFollowersConfig(userID int64) (*CursoredTweetConfigQuery, error) {
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
	item := &CursoredTweetConfigQuery{}

	if result.Item == nil {
		return item, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetFriendsConfig(userID int64) (*CursoredTweetConfigQuery, error) {
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

	item := &CursoredTweetConfigQuery{}

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

func (config *DDBDriver) GetTimelineConfig(userID int64) (*TweetConfigQuery, error) {
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

	item := &TweetConfigQuery{}

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
