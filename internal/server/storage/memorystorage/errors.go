package memorystorage

import "fmt"

type StorageGaugeNotFoundError struct {
	unknownGauge string
}

func newStorageGaugeNotFoundError(unknownGauge string) *StorageGaugeNotFoundError {
	return &StorageGaugeNotFoundError{
		unknownGauge: unknownGauge,
	}
}

func (e *StorageGaugeNotFoundError) Error() string {
	err := fmt.Sprintf("Unknown gauge: %s", e.unknownGauge)
	return err
}

type StorageCounterNotFoundError struct {
	unknownCounter string
}

func newStorageCounterNotFoundError(unknownCounter string) *StorageCounterNotFoundError {
	return &StorageCounterNotFoundError{
		unknownCounter: unknownCounter,
	}
}

func (e *StorageCounterNotFoundError) Error() string {
	err := fmt.Sprintf("Unknown Counter: %s", e.unknownCounter)
	return err
}
