package service

import "fmt"

type serviceMetricTypeUnknownError struct {
	unknownType string
}

func newServiceMetricTypeUnknownError(unknownType string) *serviceMetricTypeUnknownError {
	return &serviceMetricTypeUnknownError{
		unknownType: unknownType,
	}
}

func (e *serviceMetricTypeUnknownError) Error() string {
	err := fmt.Sprintf("Unknown metric type: %s", e.unknownType)
	return err
}

type serviceMetricNoValueUpdateError struct {
	metricName string
}

func newServiceNoValueUpdateError(metricName string) *serviceMetricNoValueUpdateError {
	return &serviceMetricNoValueUpdateError{
		metricName: metricName,
	}
}

func (e *serviceMetricNoValueUpdateError) Error() string {
	err := fmt.Sprintf("No value provided for gauge metric with name: %s", e.metricName)
	return err
}

type serviceMetricNoDeltaUpdateError struct {
	metricName string
}

func newServiceNoDeltaUpdateError(metricName string) *serviceMetricNoDeltaUpdateError {
	return &serviceMetricNoDeltaUpdateError{
		metricName: metricName,
	}
}

func (e *serviceMetricNoDeltaUpdateError) Error() string {
	err := fmt.Sprintf("No delta provided for counter metric with name: %s", e.metricName)
	return err
}

type hashNotMatchingError struct {
	requestHash string
	localHash   string
}

func newHashNotMatchingError(requestHash string, localHash string) *hashNotMatchingError {
	return &hashNotMatchingError{
		requestHash: requestHash,
		localHash:   localHash,
	}
}

func (e *hashNotMatchingError) Error() string {
	err := fmt.Sprintf("Request hash does not match with computed:\nRequest hash: %s\nComputed hash: %s", e.requestHash, e.localHash)
	return err
}
