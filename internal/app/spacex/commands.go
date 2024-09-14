package spacex

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	"spacet/pkg/logger"
)

type Commands interface {
	// SaveLaunches command persists the external launches (without a booking), in batches of 100 elements.
	// It should be called within a transaction to ensure atomicity.
	SaveExternalLaunches(ctx context.Context, launches []*domain.Launch) (err error)

	// UpdateLaunchPads reads the launchpads from the spacexapi and persists them.
	// It should be called within a transaction to ensure atomicity.
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

func (h commandsHandler) SaveExternalLaunches(ctx context.Context, launches []*domain.Launch) (err error) {
	h.l.Debug("Updating Launches")
	// batches of 100 elements
	var batchSize int = 100

	for i := 0; i < len(launches); i += batchSize {
		end := i + batchSize
		if end > len(launches) {
			end = len(launches)
		}

		batch := launches[i:end]

		if err := h.launchesRepo.SaveExternalLaunches(ctx, batch); err != nil {
			return fmt.Errorf("failed to save launches batch: %w", err)
		}
	}
	h.l.Info("Launches update completed",
		"total", len(launches))
	return nil
}

func (h commandsHandler) UpdateLaunchPads(ctx context.Context) error {
	h.l.Debug("Updating Launchpads")
	lpads, err := h.spacexClient.GetLaunchPads(ctx)
	if err != nil {
		return err
	}
	for _, lpad := range lpads {
		if _, err := h.launchPadRepo.SaveLaunchPad(ctx, lpad); err != nil {
			return err
		}
	}
	return nil
}
