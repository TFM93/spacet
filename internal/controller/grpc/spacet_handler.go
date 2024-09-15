package grpc

import (
	"context"
	gen "spacet/gen/proto/go"
	v1 "spacet/gen/proto/go"
	"spacet/internal/app"
	"spacet/internal/app/bookings"
	"spacet/internal/domain"
	"spacet/pkg/logger"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SpaceTHandler struct {
	gen.UnimplementedSpaceTServiceServer
	l                logger.Interface
	protoValidator   *protovalidate.Validator
	bookingsCommands app.BookingsServiceCommands
	bookingsQueries  app.BookingsServiceQueries
}

func (h SpaceTHandler) LaunchBooking(ctx context.Context, req *v1.BookingRequest) (*v1.Ticket, error) {
	if err := h.protoValidator.Validate(req); err != nil {
		return nil, err
	}
	ticket, err := h.bookingsCommands.BookALaunch(ctx, bookings.BookALaunchReq{
		LaunchpadID: req.GetLaunchpadId(),
		Date:        req.LaunchDate.AsTime(),
		Birthday:    req.Birthday.AsTime(),
		Destination: domain.Destination(req.DestinationId.String()),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      domain.Gender(req.GetGender().String()),
	})
	return &v1.Ticket{
		Id:            ticket.ID.String(),
		FirstName:     ticket.FirstName,
		LastName:      ticket.LastName,
		LaunchpadId:   ticket.LaunchPadID,
		LaunchDate:    timestamppb.New(ticket.LaunchDate),
		DestinationId: v1.Destination(v1.Destination_value[ticket.Destination]),
		Status:        ticket.Status,
	}, err
}

func (h SpaceTHandler) CancelBooking(ctx context.Context, req *v1.TicketID) (*v1.TicketID, error) {
	if err := h.protoValidator.Validate(req); err != nil {
		return nil, err
	}
	return req, h.bookingsCommands.CancelByID(ctx, req.Id)
}

func (h SpaceTHandler) ListBookings(ctx context.Context, req *v1.ListTicketsRequest) (*v1.ListTicketsResponse, error) {
	if req == nil {
		return nil, domain.ErrEmptyRequest
	}
	if err := h.protoValidator.Validate(req); err != nil {
		return nil, err
	}
	var resp = &gen.ListTicketsResponse{}
	ticketsList, nextCursor, err := h.bookingsQueries.ListTickets(ctx,
		bookings.ListTicketsRequest{
			Cursor: req.GetCursor(),
			Limit:  req.GetLimit(),
			Filters: domain.TicketSearchFilters{
				FirstName:   req.FirstName,
				LastName:    req.LastName,
				Destination: req.Destination,
				Status:      req.Status,
				LaunchPadID: req.LaunchpadId,
			},
		})
	if err != nil {
		return resp, err
	}
	resp.NextCursor = nextCursor
	resp.Tickets = make([]*v1.Ticket, 0, len(ticketsList))
	for _, ticket := range ticketsList {
		resp.Tickets = append(resp.Tickets, &v1.Ticket{
			Id:            ticket.ID.String(),
			FirstName:     ticket.FirstName,
			LastName:      ticket.LastName,
			LaunchDate:    timestamppb.New(ticket.LaunchDate),
			LaunchpadId:   ticket.LaunchPadID,
			DestinationId: v1.Destination(v1.Destination_value[ticket.Destination]),
			Status:        ticket.Status,
		})
	}
	return resp, nil
}
