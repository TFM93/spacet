package bookings

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	"spacet/pkg/logger"
	"time"

	"github.com/google/uuid"
)

// LaunchPadDateRestrictions map launchPadId to restricting dates
type LaunchPadDateRestrictions map[string][]time.Time

type Commands interface {
	// Cancel command cancels the launches with bookings that collide in the LaunchPadDateRestrictions field.
	Cancel(ctx context.Context, restriction LaunchPadDateRestrictions) (cancelled []uuid.UUID, err error)
	CancelByID(ctx context.Context, bookingID string) (err error)
	// BookALaunch command books a new launch if the launchpad is available.
	// Note that the booking can be cancelled if spacex needs the platform.
	BookALaunch(ctx context.Context, req BookALaunchReq) (createdBooking domain.Ticket, _ error)
}

type handler struct {
	l                logger.Interface
	bookingCmdsRepo  domain.BookingRepoCommands
	launchesCmdsRepo domain.LaunchRepoCommands
	launchesQrsRepo  domain.LaunchRepoQueries
	transaction      domain.Transaction
}

func NewCommands(logger logger.Interface, transaction domain.Transaction, bookingCmds domain.BookingRepoCommands, launchesCmds domain.LaunchRepoCommands, launchesQrs domain.LaunchRepoQueries) Commands {
	return &handler{l: logger, transaction: transaction, bookingCmdsRepo: bookingCmds, launchesCmdsRepo: launchesCmds, launchesQrsRepo: launchesQrs}
}

func (h handler) Cancel(ctx context.Context, restriction LaunchPadDateRestrictions) (cancelled []uuid.UUID, err error) {
	return h.bookingCmdsRepo.Cancel(ctx, restriction)
}

func (h handler) CancelByID(ctx context.Context, id string) (err error) {
	return h.bookingCmdsRepo.CancelByID(ctx, uuid.MustParse(id))
}

type BookALaunchReq struct {
	LaunchpadID string
	// Date will be converted into UTC
	Date        time.Time
	Destination domain.Destination
	FirstName   string
	LastName    string
	Gender      domain.Gender
	Birthday    time.Time
}

func (h handler) BookALaunch(ctx context.Context, req BookALaunchReq) (ticket domain.Ticket, _ error) {
	if !req.Destination.IsValid() {
		return ticket, domain.ErrInvalidDestination
	}

	return ticket, h.transaction.BeginTx(ctx, func(ctx context.Context) error {
		available, err := h.launchesQrsRepo.IsLaunchpadAvailableForDate(ctx, req.LaunchpadID, req.Date.UTC())
		if err != nil {
			h.l.Error(fmt.Errorf("failed to check launchpad availability: %w", err))
			return domain.ErrInternal
		}

		if !available {
			return domain.ErrLaunchPadUnavailableDate
		}

		nrLaunches, err := h.launchesQrsRepo.LaunchesOnSameDestinationOnTargetWeek(ctx, req.LaunchpadID, req.Date.UTC(), req.Destination.ToString())
		if err != nil {
			h.l.Error(fmt.Errorf("failed to check launchpad availability: %w", err))
			return domain.ErrInternal
		}

		if nrLaunches > 0 {
			return domain.ErrLaunchPadUnavailableDestination
		}

		booking := domain.Booking{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    req.Gender,
			BirthDay:  req.Birthday,
		}

		bookingID, err := h.bookingCmdsRepo.CreateBooking(ctx, booking)
		if err != nil {
			h.l.Error(fmt.Errorf("failed to create booking: %w", err))
			return domain.ErrInternal
		}

		launch := domain.Launch{
			BookingID:   bookingID,
			LaunchPadID: req.LaunchpadID,
			DateUTC:     req.Date.UTC(),
			Destination: &req.Destination,
			Domain:      domain.SpaceTDomain,
			Status:      domain.DefaultLaunchStatus,
		}

		_, err = h.launchesCmdsRepo.CreateLaunch(ctx, launch)
		if err != nil {
			h.l.Error(fmt.Errorf("failed to create launch: %w", err))
			return domain.ErrInternal
		}

		// build ticket
		ticket = domain.Ticket{
			ID:          bookingID,
			FirstName:   booking.FirstName,
			LastName:    booking.LastName,
			LaunchPadID: launch.LaunchPadID,
			LaunchDate:  launch.DateUTC,
			Status:      launch.Status,
			Destination: launch.Destination.ToString(),
		}
		return nil
	})

}
