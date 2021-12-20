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
	PutRunnerFlags(params *RunnerFlagsItem) error
	GetRunnerUsers(runnerUser *RunnerFlagsItem) ([]*RunnerFlagsItem, error)
	DeleteRunnerUser(params *RunnerFlagsItem) error
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
