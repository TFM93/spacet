package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LaunchDomain string

const (
	SpaceXDomain LaunchDomain = "SPACEX"
	SpaceTDomain LaunchDomain = "SPACET"
)

const (
	DefaultLaunchStatus = "scheduled"
)

type (

	// LaunchRepoCommands is an interface for persisting launches
	LaunchRepoCommands interface {
		// SaveLaunch creates a new launch in the database.
		// Returns the created Launch's ID and an error in case of failure.
		// If the launch already exists or a conflict is found, it returns domain.ErrLaunchAlreadyExists.
		// If an internal error occurs, it logs the error and returns domain.ErrInternal.
		CreateLaunch(ctx context.Context, launch Launch) (string, error)
		// SaveExternalLaunches creates a batch of launches in the database without associated booking.
		// If an internal error occurs, it logs the error and returns domain.ErrInternal.
		SaveExternalLaunches(ctx context.Context, launches []*Launch) error
	}

	LaunchRepoQueries interface {
		// LaunchesOnSameDestinationOnTargetWeek counts the number of launches for the date's week, launchpad and destination.
		// It expects that the implementation locks the launches for updating
		// Returns the number of launches for the date and an error in case of failure.
		// If an internal error occurs, it logs the error and returns domain.ErrInternal.
		LaunchesOnSameDestinationOnTargetWeek(ctx context.Context, launchpadID string, date time.Time, destination string) (int, error)
		// IsLaunchpadAvailableForDate checks if a given launchpad is available on the specified date.
		// It expects that the implementation locks the launches for updating
		IsLaunchpadAvailableForDate(ctx context.Context, launchpadID string, date time.Time) (bool, error)
	}

	// Launch represents a Launch in the domain model
	Launch struct {
		ID          string
		ExternalID  string
		Domain      LaunchDomain
		Name        string
		DateUTC     time.Time
		LaunchPadID string
		Destination *Destination
		Status      string
		BookingID   uuid.UUID
	}
)
