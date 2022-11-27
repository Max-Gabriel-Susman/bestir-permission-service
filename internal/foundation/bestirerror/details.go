package bestirerror

import (
	"errors"
	"sort"
)

type Detailer interface {
	error
	Details() []string
}
type detailsError struct {
	error
	details []string
}

func (d detailsError) Unwrap() error {
	return d.error
}

func (d detailsError) Details() []string {
	return d.details
}

func WithDetails(err error, de []string) error {
	return detailsError{error: err, details: de}
}

func Details(err error) []string {
	if err == nil {
		return nil
	}
	if fe := Detailer(nil); errors.As(err, &fe) {
		det := fe.Details()
		sort.Strings(det) // make details order deterministic
		return det
	}
	return nil
}
