package database

type DatabaseDriver interface {
	GetDriverName() string
	GetFavoritesConfig(userID int64) (*TweetConfigQuery, error)
	GetFollowersConfig(userID int64) (*CursoredTweetConfigQuery, error)
	GetFriendsConfig(userID int64) (*CursoredTweetConfigQuery, error)
	GetTimelineConfig(userID int64) (*TweetConfigQuery, error)
	PutFavoritesConfig(query *TweetConfigQuery) error
	PutFollowersConfig(query *CursoredTweetConfigQuery) error
	PutFriendsConfig(query *CursoredTweetConfigQuery) error
	PutTimelineConfig(query *TweetConfigQuery) error
	PutRunnerFlags(runnerName string, userid int64, flags Bits) error
	GetRunnerUsers(runner string, userID int64) ([]*RunnerFlagsItem, error)
	DeleteRunnerUser(runnerName string, userid int64) error
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
