package postgresstorage

import "fmt"

type unableToConnectToDatabaseError struct {
	connectionString string
}

func newUnableToConnectToDatabaseError(connStr string) (err *unableToConnectToDatabaseError) {
	return &unableToConnectToDatabaseError{
		connectionString: connStr,
	}
}

func (e *unableToConnectToDatabaseError) Error() string {
	err := fmt.Sprintf("Unable to connect to database: %s", e.connectionString)
	return err
}

type unableToPingDatabaseError struct {
	connectionString string
}

func newUnableToPingDatabaseError(connStr string) (err *unableToPingDatabaseError) {
	return &unableToPingDatabaseError{
		connectionString: connStr,
	}
}

func (e *unableToPingDatabaseError) Error() string {
	err := fmt.Sprintf("Unable to ping database: %s", e.connectionString)
	return err
}

type unknownTypeError struct{}

func newUnknownTypeError() (err *unknownTypeError) {
	return &unknownTypeError{}
}

func (e *unknownTypeError) Error() string {
	err := "Unknown type"
	return err
}
