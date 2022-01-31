package postgresStorage

func (st *PostgresStorage) Ping() (err error) {
	err = st.connection.Ping(st.ctx)
	if err != nil {
		return newUnableToPingDatabaseError(st.connectionString)
	}
	return nil
}
