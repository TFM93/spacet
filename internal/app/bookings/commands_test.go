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

func Test_handler_Cancel(t *testing.T) {
	tnow := time.Now()
	ctx := context.Background()
	mockedLogger := loggerMocks.NewInterface(t)
	bookingCommandsMock := domainMocks.NewBookingRepoCommands(t)
	launchesCommandsMock := domainMocks.NewLaunchRepoCommands(t)
	launchesQueriesMock := domainMocks.NewLaunchRepoQueries(t)
	transactionMock := domainMocks.NewTransaction(t)

	type args struct {
		restriction LaunchPadDateRestrictions
	}
	baseRestriction := LaunchPadDateRestrictions{
		"lpad1": []time.Time{tnow.Add(24 * time.Hour), tnow.Add(48 * time.Hour)},
	}
	baseUUID := "2f506028-4a0e-47c4-826e-5a94db8a15e7"
	tests := []struct {
		name          string
		args          args
		expectedMocks func()
		wantCancelled []uuid.UUID
		wantErr       error
	}{
		{
			name: "success",
			args: args{restriction: baseRestriction},
			expectedMocks: func() {
				bookingCommandsMock.On("Cancel", ctx, map[string][]time.Time(baseRestriction)).Return(
					[]uuid.UUID{uuid.MustParse(baseUUID)}, nil).Once()
			},
			wantCancelled: []uuid.UUID{uuid.MustParse(baseUUID)},
			wantErr:       nil,
		},
		{
			name: "withErr",
			args: args{restriction: baseRestriction},
			expectedMocks: func() {
				bookingCommandsMock.On("Cancel", ctx, map[string][]time.Time(baseRestriction)).Return(
					[]uuid.UUID{}, domain.ErrInternal).Once()
			},
			wantCancelled: []uuid.UUID{},
			wantErr:       domain.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCommands(mockedLogger, transactionMock, bookingCommandsMock, launchesCommandsMock, launchesQueriesMock)

			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}

			gotCancelled, err := h.Cancel(ctx, tt.args.restriction)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if !reflect.DeepEqual(gotCancelled, tt.wantCancelled) {
				t.Errorf("handler.Cancel() = %v, want %v", gotCancelled, tt.wantCancelled)
			}
		})
	}
}

func Test_handler_BookALaunch(t *testing.T) {
	tnow := time.Now()
	ctx := context.Background()
	mockedLogger := loggerMocks.NewInterface(t)
	bookingCommandsMock := domainMocks.NewBookingRepoCommands(t)
	launchesCommandsMock := domainMocks.NewLaunchRepoCommands(t)
	launchesQueriesMock := domainMocks.NewLaunchRepoQueries(t)
	transactionMock := domainMocks.NewTransaction(t)
	mockedLogger.On("Error", mock.Anything)
	baseUUID := "2f506028-4a0e-47c4-826e-5a94db8a15e7"
	mars := domain.DestinationMars

	tests := []struct {
		name          string
		args          BookALaunchReq
		expectedMocks func()
		wantTicket    domain.Ticket
		wantErr       error
	}{
		{
			name:       "invalid destination",
			args:       BookALaunchReq{},
			wantTicket: domain.Ticket{},
			wantErr:    domain.ErrInvalidDestination,
		}, {
			name: "failed to get launchpad availability",
			args: BookALaunchReq{
				Destination: domain.DestinationMars,
				LaunchpadID: "testID",
				Date:        tnow,
			},
			expectedMocks: func() {
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(domain.ErrInternal).Once()
				launchesQueriesMock.On("IsLaunchpadAvailableForDate", mock.Anything, "testID", tnow.UTC()).Return(false, domain.ErrInternal).Once()
			},
			wantTicket: domain.Ticket{},
			wantErr:    domain.ErrInternal,
		}, {
			name: "launchpad not available",
			args: BookALaunchReq{
				Destination: domain.DestinationMars,
				LaunchpadID: "testID",
				Date:        tnow,
			},
			expectedMocks: func() {
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(domain.ErrLaunchPadUnavailableDate).Once()
				launchesQueriesMock.On("IsLaunchpadAvailableForDate", mock.Anything, "testID", tnow.UTC()).Return(false, nil).Once()
			},
			wantTicket: domain.Ticket{},
			wantErr:    domain.ErrLaunchPadUnavailableDate,
		}, {
			name: "failed to check launchpad availability",
			args: BookALaunchReq{
				Destination: domain.DestinationMars,
				LaunchpadID: "testID",
				Date:        tnow,
			},
			expectedMocks: func() {
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(domain.ErrInternal).Once()
				launchesQueriesMock.On("IsLaunchpadAvailableForDate", mock.Anything, "testID", tnow.UTC()).Return(true, nil).Once()
				launchesQueriesMock.On("LaunchesOnSameDestinationOnTargetWeek",
					mock.Anything, "testID", tnow.UTC(),
					domain.DestinationMars.ToString()).Return(0, fmt.Errorf("something")).Once()
			},
			wantTicket: domain.Ticket{},
			wantErr:    domain.ErrInternal,
		}, {
			name: "launchpad destination unavailable",
			args: BookALaunchReq{
				Destination: domain.DestinationMars,
				LaunchpadID: "testID",
				Date:        tnow,
			},
			expectedMocks: func() {
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(domain.ErrLaunchPadUnavailableDestination).Once()
				launchesQueriesMock.On("IsLaunchpadAvailableForDate", mock.Anything, "testID", tnow.UTC()).Return(true, nil).Once()
				launchesQueriesMock.On("LaunchesOnSameDestinationOnTargetWeek",
					mock.Anything, "testID", tnow.UTC(),
					domain.DestinationMars.ToString()).Return(1, nil).Once()
			},
			wantTicket: domain.Ticket{},
			wantErr:    domain.ErrLaunchPadUnavailableDestination,
		}, {
			name: "failed to create booking",
			args: BookALaunchReq{
				Destination: domain.DestinationMars,
				LaunchpadID: "testID",
				Date:        tnow,
				FirstName:   "som",
				LastName:    "em",
				Gender:      domain.GenderFemale,
				Birthday:    tnow,
			},
			expectedMocks: func() {
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(domain.ErrInternal).Once()
				launchesQueriesMock.On("IsLaunchpadAvailableForDate", mock.Anything, "testID", tnow.UTC()).Return(true, nil).Once()
				launchesQueriesMock.On("LaunchesOnSameDestinationOnTargetWeek",
					mock.Anything, "testID", tnow.UTC(),
					domain.DestinationMars.ToString()).Return(0, nil).Once()
				bookingCommandsMock.On("CreateBooking",
					mock.Anything,
					domain.Booking{FirstName: "som", LastName: "em", Gender: domain.GenderFemale, BirthDay: tnow},
				).Return(uuid.UUID{}, fmt.Errorf("some error")).Once()
			},
			wantTicket: domain.Ticket{},
			wantErr:    domain.ErrInternal,
		}, {
			name: "failed to create launch",
			args: BookALaunchReq{
				Destination: domain.DestinationMars,
				LaunchpadID: "testID",
				Date:        tnow,
				FirstName:   "som",
				LastName:    "em",
				Gender:      domain.GenderFemale,
				Birthday:    tnow,
			},
			expectedMocks: func() {
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(domain.ErrInternal).Once()
				launchesQueriesMock.On("IsLaunchpadAvailableForDate", mock.Anything, "testID", tnow.UTC()).Return(true, nil).Once()
				launchesQueriesMock.On("LaunchesOnSameDestinationOnTargetWeek",
					mock.Anything, "testID", tnow.UTC(),
					mars.ToString()).Return(0, nil).Once()
				bookingCommandsMock.On("CreateBooking",
					mock.Anything,
					domain.Booking{FirstName: "som", LastName: "em", Gender: domain.GenderFemale, BirthDay: tnow},
				).Return(uuid.MustParse(baseUUID), nil).Once()
				launchesCommandsMock.On("CreateLaunch", mock.Anything, domain.Launch{
					Domain:      domain.SpaceTDomain,
					DateUTC:     tnow.UTC(),
					LaunchPadID: "testID",
					Destination: &mars,
					Status:      domain.DefaultLaunchStatus,
					BookingID:   uuid.MustParse(baseUUID),
				}).Return("", fmt.Errorf("something")).Once()
			},
			wantTicket: domain.Ticket{},
			wantErr:    domain.ErrInternal,
		}, {
			name: "success",
			args: BookALaunchReq{
				Destination: domain.DestinationMars,
				LaunchpadID: "testID",
				Date:        tnow,
				FirstName:   "som",
				LastName:    "em",
				Gender:      domain.GenderFemale,
				Birthday:    tnow,
			},
			expectedMocks: func() {
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(nil).Once()
				launchesQueriesMock.On("IsLaunchpadAvailableForDate", mock.Anything, "testID", tnow.UTC()).Return(true, nil).Once()
				launchesQueriesMock.On("LaunchesOnSameDestinationOnTargetWeek",
					mock.Anything, "testID", tnow.UTC(),
					mars.ToString()).Return(0, nil).Once()
				bookingCommandsMock.On("CreateBooking",
					mock.Anything,
					domain.Booking{FirstName: "som", LastName: "em", Gender: domain.GenderFemale, BirthDay: tnow},
				).Return(uuid.MustParse(baseUUID), nil).Once()
				launchesCommandsMock.On("CreateLaunch", mock.Anything, domain.Launch{
					Domain:      domain.SpaceTDomain,
					DateUTC:     tnow.UTC(),
					LaunchPadID: "testID",
					Destination: &mars,
					Status:      domain.DefaultLaunchStatus,
					BookingID:   uuid.MustParse(baseUUID),
				}).Return("something", nil).Once()
			},
			wantTicket: domain.Ticket{
				ID:          uuid.MustParse(baseUUID),
				FirstName:   "som",
				LastName:    "em",
				LaunchPadID: "testID",
				LaunchDate:  tnow.UTC(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCommands(mockedLogger, transactionMock, bookingCommandsMock, launchesCommandsMock, launchesQueriesMock)

			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}

			gotTicket, err := h.BookALaunch(ctx, tt.args)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if !reflect.DeepEqual(gotTicket, tt.wantTicket) {
				t.Errorf("handler.BookALaunch() = %v, want %v", gotTicket, tt.wantTicket)
			}
		})
	}
}
