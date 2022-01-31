package postgresStorage

import "github.com/HAGIT4/go-middle/internal/server/storage"

type PostgresStorageInterface interface {
	storage.StorageInterface

	Ping() (err error)
}
