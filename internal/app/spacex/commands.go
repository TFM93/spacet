package spacex

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	"spacet/pkg/logger"
	"time"
)

type Commands interface {
	// Handle fetches data from spaceX api and persists the information. It also checks for conflicts and cancels bookings if necessary
	// TODO: describe error handling
	Handle(ctx context.Context, req SyncSpaceXDataCommand) (err error)
	UpdateLaunchPads(ctx context.Context) (err error)
}

type SyncLaunchPadCommand struct {
	// After defines the amount of time to wait for updating again.
	// If the last_update + after is in the future, the command exits
	After time.Duration
}

type SyncSpaceXDataCommand struct {
	SyncTime time.Time
}

type dataHandler struct {
	l             logger.Interface
	launchPadRepo domain.LaunchPadRepoCommands
	launchesRepo  domain.LaunchRepoCommands
	spacexClient  domain.SpaceXAPIQueries
}

func NewCommands(logger logger.Interface, client domain.SpaceXAPIQueries, launchPadRepo domain.LaunchPadRepoCommands, launchesRepo domain.LaunchRepoCommands) Commands {
	return &dataHandler{l: logger, launchPadRepo: launchPadRepo, launchesRepo: launchesRepo, spacexClient: client}
}

func (h dataHandler) Handle(ctx context.Context, req SyncSpaceXDataCommand) (err error) {
	return fmt.Errorf("not implemented yet")
}

// UpdateLaunchPads ... todo describe.
// This function must be called within a transaction to ensure atomicity.
func (h dataHandler) UpdateLaunchPads(ctx context.Context) error {
	//todo: handle errors
	h.l.Debug("Updating Launchpads")
	lpads, err := h.spacexClient.GetLaunchPads(ctx)
	if err != nil {
		return err
	}
	for _, lpad := range lpads {
		_, err = h.launchPadRepo.SaveLaunchPad(ctx, lpad)
	}

	return err
}
