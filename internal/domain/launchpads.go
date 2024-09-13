package domain

import "context"

type (

	// LaunchPadRepoCommands is an interface for persisting launch pads
	LaunchPadRepoCommands interface {
		// SaveLaunchPad creates a new launch in the database.
		// Returns the created LaunchpPad's ID and an error in case of failure.
		// If the launchpad already exists or a conflict is found, it tries to update the launchpad.
		// If an internal error occurs, it logs the error and returns domain.ErrInternal.
		SaveLaunchPad(ctx context.Context, launch *LaunchPad) (string, error)
	}

	// LaunchPad represents a LaunchPad in the domain model
	LaunchPad struct {
		ID       string
		Name     string
		Locality string
		Region   string
		Timezone string
		Status   string
	}
)
