package api

type apiNoJSONHeaderError struct {
}

func newAPINoJSONHeaderError() *apiNoJSONHeaderError {
	return &apiNoJSONHeaderError{}
}

func (e *apiNoJSONHeaderError) Error() string {
	err := "application/json header not found in request"
	return err
}
