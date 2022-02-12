package agent

import "fmt"

type unknownSendTypeError struct{}

func newUnknownSendTypeError() *unknownSendTypeError {
	return &unknownSendTypeError{}
}

func (e *unknownSendTypeError) Error() string {
	err := "Unknown send type"
	return err
}

type unknownMetricTypeError struct {
	mType string
}

func newUnknownMetricTypeError(mType string) *unknownMetricTypeError {
	return &unknownMetricTypeError{
		mType: mType,
	}
}

func (e *unknownMetricTypeError) Error() string {
	err := fmt.Sprintf("Unknown metric type: %s", e.mType)
	return err
}
