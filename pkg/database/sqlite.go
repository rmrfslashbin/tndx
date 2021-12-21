package database

import (
	"database/sql"
	"errors"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type SqliteOption func(config *SqliteDatabaseDriver)

type SqliteDatabaseDriver struct {
	log          *logrus.Logger
	driverName   string
	databasePath string
	db           *sql.DB
}

func NewSqliteDatabase(opts ...func(*SqliteDatabaseDriver)) *SqliteDatabaseDriver {
	config := &SqliteDatabaseDriver{}
	config.driverName = "sqlite"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	sqliteDB, err := sql.Open("sqlite", config.databasePath)
	if err != nil {
		panic(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS tweets ("userid" integer NOT NULL PRIMARY KEY, "maxid" integer, "sinceid" integer, "lastupdate" integer); CREATE TABLE IF NOT EXISTS favorites ("userid" integer NOT NULL PRIMARY KEY, "maxid" integer, "sinceid" integer, "lastupdate" integer); CREATE TABLE IF NOT EXISTS followers ("userid" integer NOT NULL PRIMARY KEY, "nextCursor" integer, "previousCursor" integer, "lastupdate" integer);CREATE TABLE IF NOT EXISTS friends ("userid" integer NOT NULL PRIMARY KEY, "nextCursor" integer, "previousCursor" integer, "lastupdate" integer);`
	statement, err := sqliteDB.Prepare(createTableSQL)
	if err != nil {
		panic(err)
	}
	statement.Exec()

	config.db = sqliteDB

	return config
}

func SetDatabaseFilename(filename string) SqliteOption {
	return func(config *SqliteDatabaseDriver) {
		config.databasePath = path.Clean(filename)
	}
}

func SetSqliteLogger(logger *logrus.Logger) SqliteOption {
	return func(config *SqliteDatabaseDriver) {
		config.log = logger
	}
}

func (config *SqliteDatabaseDriver) GetDriverName() string {
	return config.driverName
}

func (config *SqliteDatabaseDriver) GetFavoritesConfig(userID int64) (*TweetConfigQuery, error) {
	var sinceID, maxID int64
	row := config.db.QueryRow(`SELECT sinceid, maxid FROM favorites WHERE userid = ?`, userID)
	err := row.Scan(&sinceID, &maxID)
	if err != nil {
		return &TweetConfigQuery{}, nil
	}
	return &TweetConfigQuery{
		UserID:  userID,
		SinceID: sinceID,
		MaxID:   maxID,
	}, err
}

func (config *SqliteDatabaseDriver) GetFollowersConfig(userID int64) (*CursoredTweetConfigQuery, error) {
	var nextCursor, previousCursor int64
	row := config.db.QueryRow(`SELECT nextCursor, previousCursor FROM followers WHERE userid = ?`, userID)
	err := row.Scan(&nextCursor, &previousCursor)
	if err != nil {
		return &CursoredTweetConfigQuery{}, nil
	}
	return &CursoredTweetConfigQuery{
		UserID:         userID,
		NextCursor:     nextCursor,
		PreviousCursor: previousCursor,
	}, err
}

func (config *SqliteDatabaseDriver) GetFriendsConfig(userID int64) (*CursoredTweetConfigQuery, error) {
	var nextCursor, previousCursor int64
	row := config.db.QueryRow(`SELECT nextCursor, previousCursor FROM friends WHERE userid = ?`, userID)
	err := row.Scan(&nextCursor, &previousCursor)
	if err != nil {
		return &CursoredTweetConfigQuery{}, nil
	}
	return &CursoredTweetConfigQuery{
		UserID:         userID,
		NextCursor:     nextCursor,
		PreviousCursor: previousCursor,
	}, err
}

func (config *SqliteDatabaseDriver) GetTimelineConfig(userID int64) (*TweetConfigQuery, error) {
	var sinceID, maxID int64
	row := config.db.QueryRow(`SELECT sinceid, maxid FROM tweets WHERE userid = ?`, userID)
	err := row.Scan(&sinceID, &maxID)
	if err != nil {
		return &TweetConfigQuery{}, nil
	}
	return &TweetConfigQuery{
		UserID:  userID,
		SinceID: sinceID,
		MaxID:   maxID,
	}, err
}

func (config *SqliteDatabaseDriver) PutFavoritesConfig(query *TweetConfigQuery) error {
	now := time.Now()
	insert := `INSERT INTO favorites (userid, sinceid, maxid, lastupdate) VALUES (?, ?, ?, ?) ON CONFLICT (userid) DO UPDATE SET sinceid = ?, maxid = ?`
	statement, err := config.db.Prepare(insert)
	if err != nil {
		return err
	}
	_, err = statement.Exec(query.UserID, query.SinceID, query.MaxID, now.UnixMilli(), query.SinceID, query.MaxID)
	config.db.Close()
	return err
}

func (config *SqliteDatabaseDriver) PutFollowersConfig(query *CursoredTweetConfigQuery) error {
	now := time.Now()
	insert := `INSERT INTO followers (userid, nextCursor, previousCursor, lastupdate) VALUES (?, ?, ?, ?) ON CONFLICT (userid) DO UPDATE SET nextCursor = ?, previousCursor = ?`
	statement, err := config.db.Prepare(insert)
	if err != nil {
		return err
	}
	_, err = statement.Exec(query.UserID, query.NextCursor, query.PreviousCursor, now.UnixMilli(), query.NextCursor, query.PreviousCursor)
	config.db.Close()
	return err
}

func (config *SqliteDatabaseDriver) PutFriendsConfig(query *CursoredTweetConfigQuery) error {
	now := time.Now()
	insert := `INSERT INTO friends (userid, nextCursor, previousCursor, lastupdate) VALUES (?, ?, ?, ?) ON CONFLICT (userid) DO UPDATE SET nextCursor = ?, previousCursor = ?`
	statement, err := config.db.Prepare(insert)
	if err != nil {
		return err
	}
	_, err = statement.Exec(query.UserID, query.NextCursor, query.PreviousCursor, now.UnixMilli(), query.NextCursor, query.PreviousCursor)
	config.db.Close()
	return err
}

func (config *SqliteDatabaseDriver) PutTimelineConfig(query *TweetConfigQuery) error {
	now := time.Now()
	createTableSQL := `INSERT INTO tweets (userid, sinceid, maxid, lastupdate) VALUES (?, ?, ?, ?) ON CONFLICT (userid) DO UPDATE SET sinceid = ?, maxid = ?`
	statement, err := config.db.Prepare(createTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(query.UserID, query.SinceID, query.MaxID, now.UnixMilli(), query.SinceID, query.MaxID)
	config.db.Close()
	return err
}

func (config *SqliteDatabaseDriver) PutRunnerFlags(params *RunnerItem) error {
	return errors.New("not implemented")
}

func (config *SqliteDatabaseDriver) GetRunnerUsers(runnerUser *RunnerItem) ([]*RunnerItem, error) {
	return nil, errors.New("not implemented")
}

func (config *SqliteDatabaseDriver) DeleteRunnerUser(params *RunnerItem) error {
	return errors.New("not implemented")
}
