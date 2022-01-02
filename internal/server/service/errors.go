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
