package storage

import "io"

type StorageDriver interface {
	Put(name string, body []byte) error
	PutStream(key string, fp io.Reader) error
	GetDriverName() string
}
