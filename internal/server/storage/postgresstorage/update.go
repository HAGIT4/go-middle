package postgresstorage

import (
	"context"

	dbModels "github.com/HAGIT4/go-middle/pkg/server/storage/models"
)

func (st *PostgresStorage) UpdateGauge(metricName string, metricValue float64) (err error) {
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()
	_, err = st.connection.Exec(ctx, "INSERT INTO gauge(id, value) VALUES($1, $2)",
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
	_, err = st.connection.Exec(ctx, "INSERT INTO counter(id, delta) VALUES($1, $2)",
		metricName, metricValue,
	)
	if err != nil {
		return err
	}
	return nil
}

func (st *PostgresStorage) UpdateBatch(req *dbModels.BatchUpdateRequest) (err error) {
	metrics := req.Metrics
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()

	tx, err := st.connection.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Prepare(ctx, "update", "INSERT INTO $1(id, $2) VALUES($3, $4)")
	if err != nil {
		return err
	}

	for _, metric := range *metrics {
		switch metric.MetricType {
		case dbModels.TypeGauge:
			_, err = tx.Exec(ctx, "update", "gauge", "value", metric.MetricID, metric.GaugeValue)
			if err != nil {
				return err
			}
		case dbModels.TypeCounter:
			_, err = tx.Exec(ctx, "update", "counter", "delta", metric.MetricID, metric.CounterDelta)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}
