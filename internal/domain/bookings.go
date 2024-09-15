package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	// BookingRepoCommands is an interface for persisting bookings
	BookingRepoCommands interface {
		// Cancel changes the status of the launches to cancelled if they meet the launchPad daily restriction (only applies to the launches with booking)
		Cancel(ctx context.Context, lPadRestriction map[string][]time.Time) (cancelled []uuid.UUID, err error)
		// CancelByID changes the status of the launch to cancelled
		CancelByID(ctx context.Context, id uuid.UUID) error
		// CreateBooking persists the provided booking and returns the autogenerated booking id
		CreateBooking(ctx context.Context, booking Booking) (id uuid.UUID, err error)
	}

	// BookingRepoQueries is an interface for query persisted bookings
	BookingRepoQueries interface {
		// ListTickets fetches a list of booked launches based on the provided filters and pagination options.
		// Parameters:
		//   cursorTicketID: ID of the booking to start listing from
		//   cursorUpdatedAt: Timestamp to filter bookings updated after this time
		//   limit: Maximum number of bookings to return
		//   filters: Criteria to filter tickets by
		// Returns a slice of booked launches objects and an error if the operation fails.
		// If the query fails to execute, it returns return domain.ErrInternal.
		// If theres an error processing the data, it returns domain.ErrFailedToProcessData.
		ListTickets(ctx context.Context, cursorTicketID string, cursorUpdatedAt *time.Time, limit int32, filters TicketSearchFilters) ([]*Ticket, error)
	}

	// Booking represents a Booking in the domain model
	Booking struct {
		ID        uuid.UUID
		FirstName string
		LastName  string
		BirthDay  time.Time
		Gender    Gender
	}

	// Ticket is the result of a successfully scheduled booking
	Ticket struct {
		ID          uuid.UUID
		FirstName   string
		LastName    string
		LaunchPadID string
		Status      string
		Destination string
		LaunchDate  time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	// TicketSearchFilters represents booked launches searchable fields
	TicketSearchFilters struct {
		FirstName   *string
		LastName    *string
		Destination *string
		Status      *string
		LaunchPadID *string
	}
)
