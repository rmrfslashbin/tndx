package database

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

type Bits uint8

type DDBOption func(config *DDBDriver)

type DDBDriver struct {
	log         *logrus.Logger
	driverName  string
	table       string
	runnerTable string
	region      string
	db          *dynamodb.Client
}

type TweetsItem struct {
	Domain     string `json:"domain"`
	UserID     int64  `json:"userid"`
	MaxID      int64  `json:"maxid"`
	SinceID    int64  `json:"sinceid"`
	LastUpdate int64  `json:"lastupdate"`
}

type FavoritesItem struct {
	Domain     string `json:"domain"`
	UserID     int64  `json:"userid"`
	MaxID      int64  `json:"maxid"`
	SinceID    int64  `json:"sinceid"`
	LastUpdate int64  `json:"lastupdate"`
}

type FollowersItem struct {
	Domain         string `json:"domain"`
	UserID         int64  `json:"userid"`
	NextCursor     int64  `json:"nextcursor"`
	PreviousCursor int64  `json:"previouscursor"`
	LastUpdate     int64  `json:"lastupdate"`
}

type FriendsItem struct {
	Domain         string `json:"domain"`
	UserID         int64  `json:"userid"`
	NextCursor     int64  `json:"nextcursor"`
	PreviousCursor int64  `json:"previouscursor"`
	LastUpdate     int64  `json:"lastupdate"`
}

type RunnerFlagsItem struct {
	RunnerName string `json:"runnername"`
	UserID     int64  `json:"userid"`
	Flags      Bits   `json:"flags"`
	LastUpdate int64  `json:"lastupdate"`
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

func SetDDBTable(table string) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.table = table
	}
}

func SetDDBRunnerTable(table string) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.runnerTable = table
	}
}

func SetDDBLogger(logger *logrus.Logger) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.log = logger
	}
}

func (config *DDBDriver) DeleteRunnerUser(params *RunnerFlagsItem) error {
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
		TableName: aws.String(config.table),
		Key: map[string]types.AttributeValue{
			"userid": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"domain": &types.AttributeValueMemberS{Value: "favorites"},
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
		TableName: aws.String(config.table),
		Key: map[string]types.AttributeValue{
			"userid": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"domain": &types.AttributeValueMemberS{Value: "followers"},
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
		TableName: aws.String(config.table),
		Key: map[string]types.AttributeValue{
			"userid": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"domain": &types.AttributeValueMemberS{Value: "friends"},
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

func (config *DDBDriver) GetRunnerUsers(runnerUsers *RunnerFlagsItem) ([]*RunnerFlagsItem, error) {
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
		config.log.Error("Error querying runner users", err)
		return nil, err
	}

	results := []*RunnerFlagsItem{}
	attributevalue.UnmarshalListOfMaps(result.Items, &results)

	return results, nil
}

func (config *DDBDriver) GetTimelineConfig(userID int64) (*TweetConfigQuery, error) {
	result, err := config.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(config.table),
		Key: map[string]types.AttributeValue{
			"userid": &types.AttributeValueMemberN{Value: strconv.FormatInt(userID, 10)},
			"domain": &types.AttributeValueMemberS{Value: "tweets"},
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
		TableName: aws.String(config.table),
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
		TableName: aws.String(config.table),
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
		TableName: aws.String(config.table),
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
		TableName: aws.String(config.table),
	}); err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutRunnerFlags(params *RunnerFlagsItem) error {
	now := time.Now()
	kvp, err := attributevalue.MarshalMap(&RunnerFlagsItem{
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
