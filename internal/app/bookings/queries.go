package bookings

import (
	"context"
	"spacet/internal/app/pagination"
	"spacet/internal/domain"
	"spacet/pkg/logger"
	"time"
)

type Queries interface {
	// ListTickets retrieves a paginated list of booked launches.
	// It supports cursor-based pagination and filtering.
	// It returns domain.ErrInvalidPaginationCursor if an invalid cursor is provided.
	// It returns domain.ErrInternal if it fails to fetch from the repository.
	ListTickets(ctx context.Context, req ListTicketsRequest) (bookings []*domain.Ticket, nextCursor string, err error)
}

type queriesHandler struct {
	l    logger.Interface
	repo domain.BookingRepoQueries
}

func NewQueries(logger logger.Interface, bookingQrsRepo domain.BookingRepoQueries) Queries {
	return &queriesHandler{l: logger, repo: bookingQrsRepo}
}

type ListTicketsRequest struct {
	Cursor  string
	Limit   int32
	Filters domain.TicketSearchFilters
}

func (h queriesHandler) ListTickets(ctx context.Context, req ListTicketsRequest) (tickets []*domain.Ticket, nextCur string, err error) {
	var updatedAtCur *time.Time
	var ticketIDCur string
	if len(req.Cursor) > 0 {
		var uCur time.Time
		uCur, ticketIDCur, err = pagination.DecodeCursor(req.Cursor)
		if err != nil {
			h.l.Debug("App-booking-queries error decoding cursor: %v", err)
			return []*domain.Ticket{}, "", domain.ErrInvalidPaginationCursor
		}
		updatedAtCur = &uCur
	}
	tickets, err = h.repo.ListTickets(ctx, ticketIDCur, updatedAtCur, req.Limit, req.Filters)
	if err != nil {
		h.l.Debug("App-booking-queries error list tickets: %v", err)
		return []*domain.Ticket{}, "", domain.ErrInternal
	}

	if len(tickets) == int(req.Limit) {
		lastTickets := tickets[len(tickets)-1]
		nextCur = pagination.EncodeCursor(lastTickets.UpdatedAt, lastTickets.ID.String())
	}

	return
}
