package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	// BookingRepoQueries is an interface for persisting bookings
	BookingRepoCommands interface {
	}

	// BookingRepoQueries is an interface for query persisted bookings
	BookingRepoQueries interface {
	}

	// Booking represents a Booking in the domain model
	// todo- fill the proper booking fields
	Booking struct {
		ID        uuid.UUID
		FirstName string
		LastName  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
