package storage

import "io"

type NoopOption func(c *LocalStorageDriver)

type NoopStorageDriver struct {
	driverName string
}

func NewNoopStorage(opts ...func(*NoopStorageDriver)) *NoopStorageDriver {
	config := &NoopStorageDriver{}
	config.driverName = "noop"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func (config *NoopStorageDriver) Put(name string, body []byte) error {
	return nil
}

func (config *NoopStorageDriver) PutStream(key string, fp io.Reader) error {
	return nil
}

func (config *NoopStorageDriver) GetDriverName() string {
	return config.driverName
}
