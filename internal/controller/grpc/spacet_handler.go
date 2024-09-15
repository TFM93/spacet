package grpc

import (
	"context"
	"fmt"
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
		Id:          ticket.ID.String(),
		FirstName:   ticket.FirstName,
		LastName:    ticket.LastName,
		LaunchpadId: ticket.LaunchPadID,
		LaunchDate:  timestamppb.New(ticket.LaunchDate),
	}, err
}

func (h SpaceTHandler) CancelBooking(ctx context.Context, req *v1.TicketID) (*v1.TicketID, error) {
	if err := h.protoValidator.Validate(req); err != nil {
		return nil, err
	}
	return req, h.bookingsCommands.CancelByID(ctx, req.Id)
}

func (h SpaceTHandler) ListBookings(ctx context.Context, req *v1.ListBookingsRequest) (*v1.ListBookingsResponse, error) {
	if err := h.protoValidator.Validate(req); err != nil {
		return nil, err
	}
	// todo-
	return nil, fmt.Errorf("not implemented yet")
}
