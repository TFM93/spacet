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

// Launches Errors
var (
	ErrLaunchAlreadyExists = fmt.Errorf("launch already exists")
)

// LaunchPad Errors
var (
	ErrLaunchPadAlreadyExists          = fmt.Errorf("launchpad already exists")
	ErrLaunchPadUnavailableDate        = fmt.Errorf("launchpad is not available on the specified date")
	ErrLaunchPadUnavailableDestination = fmt.Errorf("launchpad is not available on the specified destination on this day")
)

// Destination Errors

var (
	ErrInvalidDestination = fmt.Errorf("invalid destination")
)

// Booking Errors
var (
	ErrInvalidBookingID = fmt.Errorf("invalid booking uuid")
	ErrBookingNotFound  = fmt.Errorf("booking not found")
)
