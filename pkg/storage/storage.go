package storage

type StorageDriver interface {
	Put(name string, body []byte) error
	GetDriverName() string
}
