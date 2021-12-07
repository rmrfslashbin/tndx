package database

import "github.com/sirupsen/logrus"

type NoOpDBOption func(config *SqliteDatabaseDriver)

type NoOpDBDriver struct {
	log        *logrus.Logger
	driverName string
}

func NewNoOpDB(opts ...func(*NoOpDBDriver)) *NoOpDBDriver {
	config := &DDBDriver{}
	config.driverName = "noop"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	return config
}
