// Package search_error error handling
package search_error

import "errors"

// SearchError custom error, with stacktrace and with an error type variable
type SearchError struct {
	error
	position int
}

// New creates a new SearchError based on parameters
// message is converted into an error and stored
// typ defines the ErrType
// position is the position in the originating text where the error occurred
func New(message string, position int) *SearchError {
	err := &SearchError{
		error:    errors.New(message),
		position: position,
	}
	return err
}
