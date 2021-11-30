package service

import (
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
)

// GetUser returns a Twitter user's details.
func (config *Config) GetUser(queryParams *QueryParams) (*twitter.User, *http.Response, error) {
	// Connect to the Twitter API and fetch the requested user.
	user, resp, err := config.client.Users.Show(&twitter.UserShowParams{
		ScreenName: queryParams.ScreenName,
		UserID:     queryParams.UserID,
	})
	return user, resp, err
}

// GetUserFavorites returns a user's Twitter favorites (likes).
func (c *Config) GetUserFavorites(queryParams *QueryParams) ([]twitter.Tweet, *http.Response, error) {
	// Connect to the Twitter API and fetch the requested user's favorites.
	tweets, resp, err := c.client.Favorites.List(&twitter.FavoriteListParams{
		ScreenName: queryParams.ScreenName,
		UserID:     queryParams.UserID,
		Count:      queryParams.Count,
		SinceID:    queryParams.SinceID,
		MaxID:      queryParams.MaxID,
	})
	return tweets, resp, err
}

// GetUserFollowers returns a user's Twitter followers.
func (c *Config) GetUserFollowers(queryParams *QueryParams) (*twitter.Followers, *http.Response, error) {
	// Connect to the Twitter API and fetch timeline as defined.
	followers, resp, err := c.client.Followers.List(&twitter.FollowerListParams{
		ScreenName:          queryParams.ScreenName,
		UserID:              queryParams.UserID,
		Count:               queryParams.Count,
		SkipStatus:          queryParams.SkipStatus,
		IncludeUserEntities: queryParams.IncludeUserEntities,
		Cursor:              queryParams.Cursor,
	})
	return followers, resp, err
}

// GetUserFriends returns a user's Twitter friends.
func (c *Config) GetUserFriends(queryParams *QueryParams) (*twitter.Friends, *http.Response, error) {
	// Connect to the Twitter API and fetch timeline as defined.
	friends, resp, err := c.client.Friends.List(&twitter.FriendListParams{
		ScreenName:          queryParams.ScreenName,
		UserID:              queryParams.UserID,
		Count:               queryParams.Count,
		SkipStatus:          queryParams.SkipStatus,
		IncludeUserEntities: queryParams.IncludeUserEntities,
		Cursor:              queryParams.Cursor,
	})
	return friends, resp, err
}

// GetUserTimeline returns a user's Twitter timeline.
func (c *Config) GetUserTimeline(queryParams *QueryParams) ([]twitter.Tweet, *http.Response, error) {
	// Connect to the Twitter API and fetch timeline as defined.
	tweets, resp, err := c.client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: queryParams.ScreenName,
		UserID:     queryParams.UserID,
		Count:      queryParams.Count,
		SinceID:    queryParams.SinceID,
		MaxID:      queryParams.MaxID,
	})
	return tweets, resp, err
}
