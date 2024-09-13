package domain

import (
	"context"
	"time"
)

type LaunchDomain string

const (
	SpaceXDomain LaunchDomain = "SPACEX"
	SpaceTDomain LaunchDomain = "SPACET"
)

type (

	// LaunchRepoCommands is an interface for persisting launches
	LaunchRepoCommands interface {
		// SaveLaunch creates a new launch in the database.
		// Returns the created Launch's ID and an error in case of failure.
		// If the launch already exists or a conflict is found, it returns domain.ErrLaunchAlreadyExists.
		// If an internal error occurs, it logs the error and returns domain.ErrInternal.
		SaveLaunch(ctx context.Context, launch *Launch) (string, error)
	}

	// Launch represents a Launch in the domain model
	Launch struct {
		ID          string
		Domain      LaunchDomain
		Name        string
		DateUTC     time.Time
		DateUnix    int64
		LaunchPadID string
		Destination Destination
	}
)
