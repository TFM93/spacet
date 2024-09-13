package app

import (
	"context"
	"fmt"
	"spacet/internal/app/bookings"
	"spacet/internal/app/spacex"
	appsync "spacet/internal/app/sync"
	"spacet/internal/domain"
	"spacet/pkg/logger"
	"sync"
	"time"
)

type BookingsOrchestrator interface {
	//todo: implement
	SyncOnce(ctx context.Context, syncInterval time.Duration) error
	// StartScheduledSync starts a ticker that will trigger the SyncOnce function at each interval.
	StartScheduledSync(ctx context.Context, interval time.Duration)
	// GracefulStop stops gracefully the scheduler
	GracefulStop()
}

type handler struct {
	l                logger.Interface
	spaceXCommands   spacex.Commands
	spaceXQueries    spacex.Queries
	bookingsCommands bookings.Commands
	appsync          appsync.Commands
	stopChan         chan struct{}
	wg               sync.WaitGroup
	once             sync.Once
}

func NewBookingsOrchestrator(logger logger.Interface, spaceXCommands spacex.Commands, spaceXQueries spacex.Queries, bookingsCommands bookings.Commands, syncCommands appsync.Commands) BookingsOrchestrator {
	return &handler{
		l:                logger,
		spaceXCommands:   spaceXCommands,
		spaceXQueries:    spaceXQueries,
		bookingsCommands: bookingsCommands,
		stopChan:         make(chan struct{}),
		appsync:          syncCommands,
	}
}

func (h *handler) SyncOnce(ctx context.Context, syncInterval time.Duration) error {
	//todo: reason about returning only new launches and act on them
	return h.appsync.SyncIfNecessary(ctx, "sync_launches", syncInterval, func(ctx context.Context) error {
		// step1: get all upcoming launches
		upcoming, err := h.spaceXQueries.GetUpcomingLaunches(ctx)
		if err != nil {
			return fmt.Errorf("failed to get upcoming launches: %s", err)
		}
		launchesRestriction := make([]domain.LaunchRestriction, 0, len(upcoming))
		for _, launch := range upcoming {
			launchesRestriction = append(launchesRestriction, domain.LaunchRestriction{
				DateUTC:     launch.DateUTC,
				LaunchPadID: launch.LaunchPadID,
			})
		}
		// step2: cancel all the launches with bookings and return them
		cancelled, err := h.bookingsCommands.Cancel(ctx, launchesRestriction)
		if err != nil {
			return fmt.Errorf("failed to cancel bookings: %s", err)
		}
		h.l.Debug("%s bookings will be cancelled", len(cancelled), cancelled)
		// step3: TODO notify somehow the user by sending an event
		// step4: insert in db the upcoming launches
		if err := h.spaceXCommands.SaveLaunches(ctx, upcoming); err != nil {
			return fmt.Errorf("failed to save upcoming launches: %s", err)
		}
		return nil
	})
}

func (h *handler) StartScheduledSync(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	h.l.Info("Bookings-Scheduler: Starting")
	for {
		select {
		case <-ticker.C:
			h.wg.Add(1)
			go func() {
				defer h.wg.Done()
				if err := h.SyncOnce(ctx, interval); err != nil {
					h.l.Error("Bookings-Scheduler: %v", err)
				}
			}()
		case <-h.stopChan:
			h.l.Info("Bookings-Scheduler: Stop processing")
			ticker.Stop()
			h.wg.Wait()
			return
		case <-ctx.Done():
			h.l.Info("Bookings-Scheduler: Stop processing (context canceled)")
			ticker.Stop()
			h.wg.Wait()
			return
		}
	}
}

func (h *handler) GracefulStop() {
	h.once.Do(func() {
		close(h.stopChan)
	})
	h.wg.Wait()
}
