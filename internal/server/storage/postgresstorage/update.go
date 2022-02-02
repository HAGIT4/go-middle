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

func (st *PostgresStorage) UpdateMetric(req *dbModels.UpdateRequest) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch req.MetricType {
	case dbModels.TypeGauge:
		_, err = st.connection.Exec(ctx, "INSERT INTO gauge(id, value) VALUES($1, $2)",
			req.MetricID, req.GaugeValue,
		)
		if err != nil {
			return err
		}
	case dbModels.TypeCounter:
		_, err = st.connection.Exec(ctx, "INSERT INTO counter(id, delta) VALUES($1, $2)",
			req.MetricID, req.CounterDelta,
		)
		if err != nil {
			return err
		}
	default:
		return newUnknownTypeError()
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
	_, err = tx.Prepare(ctx, "updateGauge", "INSERT INTO gauge(id, value) VALUES($1, $2)")
	if err != nil {
		return err
	}
	_, err = tx.Prepare(ctx, "updateCounter", "INSERT INTO counter(id, delta) VALUES($1, $2)")
	if err != nil {
		return err
	}

	for _, metric := range *metrics {
		switch metric.MetricType {
		case dbModels.TypeGauge:
			_, err = tx.Exec(ctx, "updateGauge", metric.MetricID, metric.GaugeValue)
			if err != nil {
				return err
			}
		case dbModels.TypeCounter:
			_, err = tx.Exec(ctx, "updateCounter", metric.MetricID, metric.CounterDelta)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}
