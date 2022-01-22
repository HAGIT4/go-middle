package agent

type unknownSendTypeError struct{}

func newUnknownSendTypeError() *unknownSendTypeError {
	return &unknownSendTypeError{}
}

func (e *unknownSendTypeError) Error() string {
	err := "Unknown send type"
	return err
}
