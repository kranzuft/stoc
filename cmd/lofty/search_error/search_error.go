package search_error

import "errors"

// SearchError custom error, with stacktrace and with an error type variable
type SearchError struct {
	error
	typ      ErrType
	position int
}

type ErrType int

const (
	MissingRightBracket ErrType = iota
	MismatchedBrackets
)

func New(message string, typ ErrType, position int) *SearchError {
	err := &SearchError{
		error:    errors.New(message),
		typ:      typ,
		position: position,
	}
	return err
}
