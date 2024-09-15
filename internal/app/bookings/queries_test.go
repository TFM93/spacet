package bookings

import (
	"context"
	"fmt"
	"reflect"
	domainMocks "spacet/gen/mocks/spacet/domain"
	loggerMocks "spacet/gen/mocks/spacet/pkg/logger"
	"spacet/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_queriesHandler_ListTickets(t *testing.T) {
	bookingQrsMock := domainMocks.NewBookingRepoQueries(t)
	mockedLogger := loggerMocks.NewInterface(t)
	ctx := context.Background()
	tnow := time.Now()

	domainTicket1 := domain.Ticket{
		ID:          uuid.MustParse("c12e23f3-f5e3-41bc-aeca-9d66bd0b96a3"),
		FirstName:   "nick",
		LastName:    "nock",
		LaunchPadID: "123",
		Status:      "scheduled",
		Destination: domain.DestinationEuropa.ToString(),
		LaunchDate:  tnow.Add(24 * time.Hour),
		CreatedAt:   tnow,
		UpdatedAt:   tnow,
	}
	domainTicket2 := domain.Ticket{
		ID:          uuid.MustParse("c12e23f3-f5e3-41bc-f5e3-9d66bd0b96a3"),
		FirstName:   "nick2",
		LastName:    "nock2",
		LaunchPadID: "456",
		Status:      "scheduled",
		Destination: domain.DestinationAsteroidBelt.ToString(),
		LaunchDate:  tnow.Add(48 * time.Hour),
		CreatedAt:   tnow,
		UpdatedAt:   tnow,
	}

	expectedCursorTime := time.Date(2024, time.August, 22, 20, 9, 11, 938220000, time.Local)
	var nilt *time.Time

	tests := []struct {
		name          string
		req           ListTicketsRequest
		wantTickets   []*domain.Ticket
		expectedMocks func()
		wantNextCur   bool
		wantErr       error
	}{
		{
			name: "error listing users",
			req: ListTicketsRequest{
				Cursor: "",
				Limit:  2,
			},
			expectedMocks: func() {
				mockedLogger.On("Debug", mock.Anything, mock.Anything).Once()
				bookingQrsMock.On("ListTickets", mock.Anything, "", nilt, int32(2), domain.TicketSearchFilters{}).Return(
					[]*domain.Ticket{}, fmt.Errorf("something went wrong")).Once()
			},
			wantTickets: []*domain.Ticket{},
			wantNextCur: false,
			wantErr:     domain.ErrInternal,
		},
		{
			name: "more pages",
			req: ListTicketsRequest{
				Cursor: "",
				Limit:  2,
			},
			expectedMocks: func() {
				bookingQrsMock.On("ListTickets", mock.Anything, "", nilt, int32(2), domain.TicketSearchFilters{}).Return(
					[]*domain.Ticket{&domainTicket1, &domainTicket2}, nil).Once()
			},
			wantTickets: []*domain.Ticket{&domainTicket1, &domainTicket2},
			wantNextCur: true,
			wantErr:     nil,
		},
		{
			name: "with cursor",
			req: ListTicketsRequest{
				Cursor: "MjAyNC0wOC0yMlQyMDowOToxMS45MzgyMiswMTowMHxjMTJlMjNmMy1mNWUzLTQxYmMtYWVjYS05ZDY2YmQwYjk2YTM=",
				Limit:  2,
			},
			expectedMocks: func() {
				bookingQrsMock.On("ListTickets", mock.Anything, "c12e23f3-f5e3-41bc-aeca-9d66bd0b96a3", &expectedCursorTime, int32(2), domain.TicketSearchFilters{}).Return(
					[]*domain.Ticket{&domainTicket1, &domainTicket2}, nil).Once()
			},
			wantTickets: []*domain.Ticket{&domainTicket1, &domainTicket2},
			wantNextCur: true,
			wantErr:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewQueries(mockedLogger, bookingQrsMock)

			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}

			gotTickets, gotNextCur, err := h.ListTickets(ctx, tt.req)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if !reflect.DeepEqual(gotTickets, tt.wantTickets) {
				t.Errorf("queriesHandler.ListTickets() gotTickets = %v, want %v", gotTickets, tt.wantTickets)
			}
			if (len(gotNextCur) > 0) != tt.wantNextCur {
				t.Errorf("queriesHandler.ListTickets() gotNextCur = %v, want %v", gotNextCur, tt.wantNextCur)
			}
		})
	}
}
