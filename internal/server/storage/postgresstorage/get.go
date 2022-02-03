package postgresstorage

import (
	"context"
	"fmt"
)

func (st *PostgresStorage) GetGauge(metricName string) (metricValue float64, err error) {
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()
	sqlResult, err := st.connection.Query(ctx, "SELECT value FROM gauge WHERE id=$1", metricName)
	if err != nil {
		return 0, err
	}
	defer sqlResult.Close()

	for sqlResult.Next() {
		if err = sqlResult.Scan(&metricValue); err != nil {
			return 0, err
		}
	}
	err = sqlResult.Err()
	if err != nil {
		return 0, err
	}
	fmt.Println("Got gauge from db:", metricValue)
	return metricValue, nil
}

func (st *PostgresStorage) GetGaugeAll() (metricNameToValue map[string]float64, err error) {
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()
	sqlResult, err := st.connection.Query(ctx, "SELECT * FROM gauge")
	if err != nil {
		return nil, err
	}
	defer sqlResult.Close()

	var metric string
	var value float64
	metricNameToValue = make(map[string]float64)
	for sqlResult.Next() {
		if err = sqlResult.Scan(&metric, &value); err != nil {
			return nil, err
		}
		metricNameToValue[metric] = value
	}
	err = sqlResult.Err()
	if err != nil {
		return nil, err
	}
	return metricNameToValue, nil
}

func (st *PostgresStorage) GetCounter(metricName string) (metricValue int64, err error) {
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()
	sqlResult, err := st.connection.Query(ctx, "SELECT delta FROM counter WHERE id=$1", metricName)
	if err != nil {
		return 0, err
	}
	defer sqlResult.Close()

	for sqlResult.Next() {
		if err = sqlResult.Scan(&metricValue); err != nil {
			return 0, err
		}
	}
	fmt.Println("Got counter from db:", metricValue)
	return metricValue, nil
}

func (st *PostgresStorage) GetCounterAll() (metricNameToValue map[string]int64, err error) {
	ctx, cancel := context.WithCancel(st.ctx)
	defer cancel()
	sqlResult, err := st.connection.Query(ctx, "SELECT * FROM counter")
	if err != nil {
		return nil, err
	}
	defer sqlResult.Close()

	var metric string
	var value int64
	metricNameToValue = make(map[string]int64)
	for sqlResult.Next() {
		if err = sqlResult.Scan(&metric, &value); err != nil {
			return nil, err
		}
		metricNameToValue[metric] = value
	}
	err = sqlResult.Err()
	if err != nil {
		return nil, err
	}
	return metricNameToValue, nil
}
