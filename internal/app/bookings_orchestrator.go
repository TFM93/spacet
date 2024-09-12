package app

import (
	"context"
	"fmt"
	"spacet/internal/app/bookings"
	"spacet/internal/app/spacex"
	"spacet/pkg/logger"
	"sync"
	"time"
)

type BookingsOrchestrator interface {
	//todo: implement
	SyncOnce(ctx context.Context) error
	// StartScheduledSync starts a ticker that will trigger the SyncOnce function at each interval.
	StartScheduledSync(ctx context.Context, interval time.Duration)
	// GracefulStop stops gracefully the scheduler
	GracefulStop()
}

type handler struct {
	l                logger.Interface
	spaceXCommands   spacex.Commands
	bookingsCommands bookings.Commands
	stopChan         chan struct{}
	wg               sync.WaitGroup
	once             sync.Once
}

func NewBookingsOrchestrator(logger logger.Interface, spaceXCommands spacex.Commands, bookingsCommands bookings.Commands) BookingsOrchestrator {
	return &handler{
		l:                logger,
		spaceXCommands:   spaceXCommands,
		bookingsCommands: bookingsCommands,
		stopChan:         make(chan struct{}),
	}
}

func (h *handler) SyncOnce(ctx context.Context) error {
	//todo: sync spaceX data, return upcoming launches and cancel bookings if necessary
	//todo: reason about returning only new launches and act on them
	return fmt.Errorf("not implemented yet")
}

func (h *handler) StartScheduledSync(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h.wg.Add(1)
			go func() {
				defer h.wg.Done()
				if err := h.SyncOnce(ctx); err != nil {
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
	h.once.Do(func() { close(h.stopChan) })
	h.wg.Wait()
}
