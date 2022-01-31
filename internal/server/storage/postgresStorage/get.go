package postgresstorage

func (st *PostgresStorage) GetGauge(metricName string) (metricValue float64, err error) {
	return 0, nil
}

func (st *PostgresStorage) GetGaugeAll() (metricNameToValue map[string]float64, err error) {
	return nil, nil
}

func (st *PostgresStorage) GetCounter(metricName string) (metricValue int64, err error) {
	return 0, nil
}

func (st *PostgresStorage) GetCounterAll() (metricNameToValue map[string]int64, err error) {
	return nil, nil
}
