// Package postgresstorage postgres implementation of storage
package postgresstorage

import (
	"context"

	"github.com/HAGIT4/go-middle/pkg/server/storage/config"
	"github.com/jackc/pgx/v4"
)

type PostgresStorage struct {
	connectionString string

	ctx        context.Context
	connection *pgx.Conn
}

var _ PostgresStorageInterface = (*PostgresStorage)(nil)

func NewPostgresStorage(cfg *config.PostgresStorageConfig) (st *PostgresStorage, err error) {
	ctx := context.Background()
	ctxT, cancel := context.WithCancel(ctx)
	defer cancel()

	connCfg, err := pgx.ParseConfig(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctxT, connCfg)
	if err != nil {
		return nil, newUnableToConnectToDatabaseError(cfg.ConnectionString)
	}

	_, err = conn.Exec(ctxT, "CREATE TABLE IF NOT EXISTS gauge (id TEXT, value double precision)")
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(ctxT, "CREATE TABLE IF NOT EXISTS counter (id TEXT, delta bigint)")
	if err != nil {
		return nil, err
	}

	st = &PostgresStorage{
		connectionString: cfg.ConnectionString,
		ctx:              ctx,
		connection:       conn,
	}
	return st, nil
}
