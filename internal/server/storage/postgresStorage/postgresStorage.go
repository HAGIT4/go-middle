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

func (st *PostgresStorage) Ping() (err error) {
	err = st.connection.Ping(st.ctx)
	if err != nil {
		return newUnableToPingDatabaseError(st.connectionString)
	}
	return nil
}
