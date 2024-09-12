package domain

import "fmt"

// Generic Errors
var (
	ErrInternal                = fmt.Errorf("internal error")
	ErrInvalidPW               = fmt.Errorf("invalid password")
	ErrFailedToProcessData     = fmt.Errorf("failed to process data")
	ErrInvalidPaginationCursor = fmt.Errorf("cursor must be a base64 string")
	ErrEmptyRequest            = fmt.Errorf("empty request")
)
