package api

import "fmt"

type HandlerForbiddenMethodError struct {
	forbiddenMethod string
}

func newHandlerForbiddenMethodError(forbiddenMethod string) *HandlerForbiddenMethodError {
	return &HandlerForbiddenMethodError{
		forbiddenMethod: forbiddenMethod,
	}
}

func (e *HandlerForbiddenMethodError) Error() string {
	err := fmt.Sprintf("Forbidden method: %s", e.forbiddenMethod)
	return err
}

type HandlerApplicationTypeError struct {
	invalidType string
}

func newHandlerApplicationTypeError(invalidType string) *HandlerApplicationTypeError {
	return &HandlerApplicationTypeError{
		invalidType: invalidType,
	}
}

func (e *HandlerApplicationTypeError) Error() string {
	err := fmt.Sprintf("Invalid application type header: %s", e.invalidType)
	return err
}

type UpdateUnknownMetricTypeError struct {
	unknownMetricType string
}

func newUpdateUnknownMetricTypeError(unknownMetricType string) *UpdateUnknownMetricTypeError {
	return &UpdateUnknownMetricTypeError{
		unknownMetricType: unknownMetricType,
	}
}

func (e *UpdateUnknownMetricTypeError) Error() string {
	err := fmt.Sprintf("Unknown metric type: %s", e.unknownMetricType)
	return err
}

type UpdateInvalidPathFormatError struct{}

func newUpdateInvalidPathFormatError() *UpdateInvalidPathFormatError {
	return &UpdateInvalidPathFormatError{}
}

func (e *UpdateInvalidPathFormatError) Error() string {
	err := "Invalid path format"
	return err
}
