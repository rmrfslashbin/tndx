package database

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type Bits uint8

type DDBOption func(config *SqliteDatabaseDriver)

type DDBDriver struct {
	log         *logrus.Logger
	driverName  string
	table       string
	runnerTable string
	region      string
	db          *dynamodb.DynamoDB
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
)

func NewDDB(opts ...func(*DDBDriver)) *DDBDriver {
	config := &DDBDriver{}
	config.driverName = "ddb"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	config.db = svc

	return config
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

func (config *DDBDriver) GetDriverName() string {
	return config.driverName
}

func (config *DDBDriver) GetFavoritesConfig(userID int64) (*TweetConfigQuery, error) {
	result, err := config.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(config.table),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				N: aws.String(strconv.FormatInt(userID, 10)),
			},
			"domain": {
				S: aws.String("favorites"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	item := &TweetConfigQuery{}

	if result.Item == nil {
		return item, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetFollowersConfig(userID int64) (*CursoredTweetConfigQuery, error) {
	result, err := config.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(config.table),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				N: aws.String(strconv.FormatInt(userID, 10)),
			},
			"domain": {
				S: aws.String("followers"),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	item := &CursoredTweetConfigQuery{}

	if result.Item == nil {
		return item, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetFriendsConfig(userID int64) (*CursoredTweetConfigQuery, error) {
	result, err := config.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(config.table),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				N: aws.String(strconv.FormatInt(userID, 10)),
			},
			"domain": {
				S: aws.String("friends"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	item := &CursoredTweetConfigQuery{}

	if result.Item == nil {
		return item, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) GetTimelineConfig(userID int64) (*TweetConfigQuery, error) {
	result, err := config.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(config.table),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				N: aws.String(strconv.FormatInt(userID, 10)),
			},
			"domain": {
				S: aws.String("tweets"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	item := &TweetConfigQuery{}

	if result.Item == nil {
		return item, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *DDBDriver) PutFavoritesConfig(query *TweetConfigQuery) error {
	now := time.Now()
	item := &FavoritesItem{
		Domain:     "favorites",
		UserID:     query.UserID,
		MaxID:      query.MaxID,
		SinceID:    query.SinceID,
		LastUpdate: now.UnixMilli(),
	}
	kvp, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.table),
	}

	_, err = config.db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutFollowersConfig(query *CursoredTweetConfigQuery) error {
	now := time.Now()
	item := &FollowersItem{
		Domain:         "followers",
		UserID:         query.UserID,
		PreviousCursor: query.PreviousCursor,
		NextCursor:     query.NextCursor,
		LastUpdate:     now.UnixMilli(),
	}
	kvp, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.table),
	}

	_, err = config.db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutFriendsConfig(query *CursoredTweetConfigQuery) error {
	now := time.Now()
	item := &FriendsItem{
		Domain:         "friends",
		UserID:         query.UserID,
		PreviousCursor: query.PreviousCursor,
		NextCursor:     query.NextCursor,
		LastUpdate:     now.UnixMilli(),
	}
	kvp, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.table),
	}

	_, err = config.db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutTimelineConfig(query *TweetConfigQuery) error {
	now := time.Now()
	item := &TweetsItem{
		Domain:     "tweets",
		UserID:     query.UserID,
		MaxID:      query.MaxID,
		SinceID:    query.SinceID,
		LastUpdate: now.UnixMilli(),
	}
	kvp, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.table),
	}

	_, err = config.db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) PutRunnerFlags(runnerName string, userid int64, flags Bits) error {
	now := time.Now()
	item := &RunnerFlagsItem{
		RunnerName: runnerName,
		UserID:     userid,
		Flags:      flags,
		LastUpdate: now.UnixMilli(),
	}
	kvp, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      kvp,
		TableName: aws.String(config.runnerTable),
	}

	_, err = config.db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (config *DDBDriver) GetRunnerUsers(runner string, userID int64) ([]*RunnerFlagsItem, error) {
	var input *dynamodb.QueryInput
	if userID == 0 {
		input = &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v1": {
					S: aws.String(runner),
				},
			},
			KeyConditionExpression: aws.String("runnername = :v1"),
			TableName:              aws.String(config.runnerTable),
		}
	} else {
		input = &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v1": {
					S: aws.String(runner),
				},
				":userid": {
					N: aws.String(strconv.FormatInt(userID, 10)),
				},
			},
			KeyConditionExpression: aws.String("runnername = :v1 and userid = :userid"),
			TableName:              aws.String(config.runnerTable),
		}
	}

	result, err := config.db.Query(input)
	if err != nil {
		return nil, err
	}

	results := []*RunnerFlagsItem{}
	for _, i := range result.Items {
		item := &RunnerFlagsItem{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (config *DDBDriver) DeleteRunnerUser(runnerName string, userid int64) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"runnername": {
				S: aws.String(runnerName),
			},
			"userid": {
				N: aws.String(strconv.FormatInt(userid, 10)),
			},
		},
		TableName: aws.String(config.runnerTable),
	}
	_, err := config.db.DeleteItem(input)
	return err
}
