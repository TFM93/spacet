package app

import (
	"context"
	"encoding/json"
	"fmt"
	"spacet/internal/domain"
	"spacet/pkg/logger"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BookingsOrchestrator interface {
	//SyncOnce syncs the spacexapi launches and ensures consistency with the internal launches (with associated booking)
	SyncOnce(ctx context.Context, syncInterval time.Duration) error
	// StartScheduledSync starts a ticker that will trigger the SyncOnce function at each interval.
	StartScheduledSync(ctx context.Context, interval time.Duration)
	// GracefulStop stops gracefully the scheduler
	GracefulStop()
}

type handler struct {
	l                logger.Interface
	spaceXCommands   SpaceXServiceCommands
	spaceXQueries    SpaceXServiceQueries
	bookingsCommands BookingsServiceCommands
	appsync          SyncServiceCommands
	outboxRepo       domain.OutboxRepoCommands
	stopChan         chan struct{}
	wg               sync.WaitGroup
	once             sync.Once
}

func NewBookingsOrchestrator(logger logger.Interface,
	spaceXCommands SpaceXServiceCommands,
	spaceXQueries SpaceXServiceQueries,
	bookingsCommands BookingsServiceCommands,
	syncCommands SyncServiceCommands,
	outboxRepo domain.OutboxRepoCommands) BookingsOrchestrator {
	return &handler{
		l:                logger,
		spaceXCommands:   spaceXCommands,
		spaceXQueries:    spaceXQueries,
		bookingsCommands: bookingsCommands,
		stopChan:         make(chan struct{}),
		appsync:          syncCommands,
		outboxRepo:       outboxRepo,
	}
}

func (h *handler) sendCancelledNotification(txCtx context.Context, cancelled []uuid.UUID) error {
	if len(cancelled) == 0 {
		return nil
	}
	pl := domain.BookingsCancelledEventPayload{}
	pl.FromUUIDs(cancelled)

	payload, err := json.Marshal(pl)
	if err != nil {
		return err
	}
	event := &domain.Event{
		Type:    "BookingsCancelled",
		Payload: payload,
	}
	if _, err := h.outboxRepo.AddEvent(txCtx, event); err != nil {
		return err
	}
	return nil
}

func (h *handler) SyncOnce(ctx context.Context, syncInterval time.Duration) error {
	return h.appsync.SyncIfNecessary(ctx, "sync_launches", syncInterval, func(txCtx context.Context) error {
		upcoming, err := h.spaceXQueries.GetUpcomingLaunches(txCtx)
		if err != nil {
			return fmt.Errorf("failed to get upcoming launches: %s", err)
		}
		datesPerLPad := make(map[string][]time.Time, len(upcoming))
		for _, launch := range upcoming {
			datesPerLPad[launch.LaunchPadID] = append(datesPerLPad[launch.LaunchPadID], launch.DateUTC)
		}

		// cancel all the launches with bookings and return them
		cancelled, err := h.bookingsCommands.Cancel(txCtx, datesPerLPad)
		if err != nil {
			return fmt.Errorf("failed to cancel bookings: %s", err)
		}
		h.l.Debug("%s bookings will be cancelled", len(cancelled), cancelled)

		if err := h.spaceXCommands.SaveExternalLaunches(txCtx, upcoming); err != nil {
			return fmt.Errorf("failed to save upcoming launches: %s", err)
		}

		return h.sendCancelledNotification(txCtx, cancelled)
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
