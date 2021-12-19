package storage

import (
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
)

type LocalOption func(c *LocalStorageDriver)

type LocalStorageDriver struct {
	driverName string
	rootPath   string
	dirMode    fs.FileMode
	fileMode   fs.FileMode
}

func NewLocalStorage(opts ...func(*LocalStorageDriver)) *LocalStorageDriver {
	config := &LocalStorageDriver{}
	config.driverName = "local"
	config.dirMode = 0755
	config.fileMode = 0644

	cwd, err := os.Getwd()
	if err != nil {
		config.rootPath = "/tmp"
	} else {
		config.rootPath = path.Join((cwd), "data")
	}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if err := os.MkdirAll(config.rootPath, config.dirMode); err != nil {
		panic(err)
	}
	return config
}

func SetDirMode(dirMode os.FileMode) LocalOption {
	return func(config *LocalStorageDriver) {
		config.dirMode = dirMode
	}
}

func SetFileMode(fileMode os.FileMode) LocalOption {
	return func(config *LocalStorageDriver) {
		config.fileMode = fileMode
	}
}

func SetRootPath(rootPath string) LocalOption {
	return func(config *LocalStorageDriver) {
		config.rootPath = path.Clean(rootPath)
	}
}

func (config *LocalStorageDriver) Put(key string, body []byte) error {
	// Append ".gz" to the key (filename).
	key = key + ".gz"

	// Create a file to write to.
	fqpn := path.Join(config.rootPath, key)

	if err := os.MkdirAll(path.Dir(fqpn), config.dirMode); err != nil {
		return err
	}

	file, err := os.Create(fqpn)
	if err != nil {
		return err
	}
	defer file.Close()

	// gzip data
	zw := gzip.NewWriter(file)
	_, err = zw.Write(body)

	// Bail out if we got an error while compressing.
	if err := zw.Close(); err != nil {
		return err
	}

	return err
}

func (config *LocalStorageDriver) PutStream(key string, fp io.Reader) error {
	return errors.New("not implemented")
}

func (config *LocalStorageDriver) GetDriverName() string {
	return config.driverName
}
