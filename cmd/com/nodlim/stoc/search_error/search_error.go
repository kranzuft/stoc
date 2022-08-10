// Package search_error error handling
package search_error

import "errors"

// SearchError custom error
type SearchError struct {
	error
	position int
}

// getPos getter for position in SearchError
func (err SearchError) getPos() int {
	return err.position
}

// PosError custom error with position interface
type PosError interface {
	error
	getPos() int
}

// New creates a new SearchError based on parameters
// message is converted into an error and stored
// typ defines the ErrType
// Position is the Position in the originating text where the error occurred
func New(message string, position int) PosError {
	err := &SearchError{
		error:    errors.New(message),
		position: position,
	}
	return err
}
