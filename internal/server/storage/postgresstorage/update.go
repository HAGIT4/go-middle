package postgresstorage

import "context"

func (st *PostgresStorage) UpdateGauge(metricName string, metricValue float64) (err error) {
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()
	_, err = st.connection.Exec(ctx, "INSERT INTO gauge (id, value) VALUES ($1, $2)",
		metricName, metricValue,
	)
	if err != nil {
		return err
	}
	return nil
}

func (st *PostgresStorage) UpdateCounter(metricName string, metricValue int64) (err error) {
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()
	_, err = st.connection.Exec(ctx, "INSERT INTO counter (id, delta) VALUES ($1, $2)",
		metricName, metricValue,
	)
	if err != nil {
		return err
	}
	return nil
}
