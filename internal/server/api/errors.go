package api

type apiNoJSONHeaderError struct {
}

func newApiNoJSONHeaderError() *apiNoJSONHeaderError {
	return &apiNoJSONHeaderError{}
}

func (e *apiNoJSONHeaderError) Error() string {
	err := "application/json header not found in request"
	return err
}
