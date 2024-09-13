package spacex

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	"spacet/pkg/logger"
)

type Commands interface {
	// TODO: describe
	SaveLaunches(ctx context.Context, launches []*domain.Launch) (err error)
	UpdateLaunchPads(ctx context.Context) (err error)
}

type commandsHandler struct {
	l             logger.Interface
	launchPadRepo domain.LaunchPadRepoCommands
	launchesRepo  domain.LaunchRepoCommands
	spacexClient  domain.SpaceXAPIQueries
}

func NewCommands(logger logger.Interface, client domain.SpaceXAPIQueries, launchPadRepo domain.LaunchPadRepoCommands, launchesRepo domain.LaunchRepoCommands) Commands {
	return &commandsHandler{l: logger, launchPadRepo: launchPadRepo, launchesRepo: launchesRepo, spacexClient: client}
}

// UpdateLaunches ... todo describe.
// This function must be called within a transaction to ensure atomicity.
func (h commandsHandler) SaveLaunches(ctx context.Context, launches []*domain.Launch) (err error) {
	h.l.Debug("Updating Launches")
	// batches of 100 elements
	var batchSize int = 100

	for i := 0; i < len(launches); i += batchSize {
		end := i + batchSize
		if end > len(launches) {
			end = len(launches)
		}

		batch := launches[i:end]

		if err := h.launchesRepo.SaveLaunchesBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to save launches batch: %w", err)
		}
	}
	h.l.Info("Launches update completed",
		"total", len(launches))
	return nil
}

// UpdateLaunchPads ... todo describe.
// This function must be called within a transaction to ensure atomicity.
func (h commandsHandler) UpdateLaunchPads(ctx context.Context) error {
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
