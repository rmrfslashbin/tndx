package database

type DatabaseDriver interface {
	DeleteRunnerUser(*RunnerItem) error
	GetDriverName() string
	GetFavoritesConfig(int64) (*TweetConfigQuery, error)
	GetFollowersConfig(int64) (*CursoredTweetConfigQuery, error)
	GetFriendsConfig(int64) (*CursoredTweetConfigQuery, error)
	GetRunnerUsers(*RunnerItem) ([]*RunnerItem, error)
	GetTimelineConfig(int64) (*TweetConfigQuery, error)
	PutFavoritesConfig(*TweetConfigQuery) error
	PutFollowersConfig(*CursoredTweetConfigQuery) error
	PutFriendsConfig(*CursoredTweetConfigQuery) error
	PutTimelineConfig(*TweetConfigQuery) error
	PutRunnerFlags(*RunnerItem) error
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
