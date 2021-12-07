package database

import "github.com/sirupsen/logrus"

type NoopDBOption func(config *NoopDBDriver)

type NoopDBDriver struct {
	log        *logrus.Logger
	driverName string
}

func NewNoopDB(opts ...func(*NoopDBDriver)) *NoopDBDriver {
	config := &NoopDBDriver{}
	config.driverName = "noop"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func SetNoopDBLogger(log *logrus.Logger) NoopDBOption {
	return func(config *NoopDBDriver) {
		config.log = log
	}
}

func (config *NoopDBDriver) GetDriverName() string {
	return config.driverName
}

func (config *NoopDBDriver) GetFavoritesConfig(userID int64) (*TweetConfigQuery, error) {
	return nil, nil
}

func (config *NoopDBDriver) GetFollowersConfig(userID int64) (*CursoredTweetConfigQuery, error) {
	return nil, nil
}

func (config *NoopDBDriver) GetFriendsConfig(userID int64) (*CursoredTweetConfigQuery, error) {
	return nil, nil
}

func (config *NoopDBDriver) GetTimelineConfig(userID int64) (*TweetConfigQuery, error) {
	return nil, nil
}

func (config *NoopDBDriver) PutFavoritesConfig(query *TweetConfigQuery) error {
	return nil
}

func (config *NoopDBDriver) PutFollowersConfig(query *CursoredTweetConfigQuery) error {
	return nil
}

func (config *NoopDBDriver) PutFriendsConfig(query *CursoredTweetConfigQuery) error {
	return nil
}

func (config *NoopDBDriver) PutTimelineConfig(query *TweetConfigQuery) error {
	return nil
}

func (config *NoopDBDriver) PutRunnerFlags(runnerName string, userid int64, flags Bits) error {
	return nil
}

func (config *NoopDBDriver) GetRunnerUsers(runner string, userID int64) ([]*RunnerFlagsItem, error) {
	return nil, nil
}

func (config *NoopDBDriver) DeleteRunnerUser(runnerName string, userid int64) error {
	return nil
}