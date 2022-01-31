package postgresStorage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type PostgresStorage struct {
	connectionString string

	ctx        context.Context
	connection *pgx.Conn
}

var _ PostgresStorageInterface = (*PostgresStorage)(nil)

func NewPostgresStorage(cfg *PostgresStorageConfig) (st *PostgresStorage, err error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.ConnectionString)
	if err != nil {
		return nil, newUnableToConnectToDatabaseError(cfg.ConnectionString)
	}
	st = &PostgresStorage{
		connectionString: cfg.ConnectionString,
		ctx:              ctx,
		connection:       conn,
	}
	return st, nil
}
