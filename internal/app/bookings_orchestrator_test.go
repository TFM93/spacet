package app

import (
	"context"
	"fmt"
	serviceMocks "spacet/gen/mocks/spacet/app"
	loggerMocks "spacet/gen/mocks/spacet/pkg/logger"
	"spacet/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_handler_SyncOnce(t *testing.T) {
	tnow := time.Now()
	ctx := context.Background()
	mockedLogger := loggerMocks.NewInterface(t)
	bookingsCmdsMock := serviceMocks.NewBookingsServiceCommands(t)
	syncCmdsMock := serviceMocks.NewSyncServiceCommands(t)
	spaceXCmdsMock := serviceMocks.NewSpaceXServiceCommands(t)
	spaceXQrsMock := serviceMocks.NewSpaceXServiceQueries(t)

	upcomingLaunches := []*domain.Launch{
		{LaunchPadID: "test", DateUTC: tnow.UTC()},
		{LaunchPadID: "test2", DateUTC: tnow.UTC().Add(24 * time.Hour)},
	}
	tests := []struct {
		name          string
		expectedMocks func()
		wantErr       error
	}{
		{
			name:    "failed to get upcoming launches",
			wantErr: fmt.Errorf("failed to get upcoming launches: something"),
			expectedMocks: func() {
				syncCmdsMock.On("SyncIfNecessary",
					mock.Anything, "sync_launches", 0*time.Hour,
					mock.AnythingOfType("domain.SyncAction")).Run(func(args mock.Arguments) {
					fn := args.Get(3).(domain.SyncAction)
					fn(args.Get(0).(context.Context))
				}).Return(fmt.Errorf("failed to get upcoming launches: something")).Once()
				spaceXQrsMock.On("GetUpcomingLaunches", mock.Anything).Return([]*domain.Launch{}, fmt.Errorf("something")).Once()
			},
		}, {
			name:    "failed to cancel conflicting bookings",
			wantErr: fmt.Errorf("failed to cancel bookings: something"),
			expectedMocks: func() {
				syncCmdsMock.On("SyncIfNecessary",
					mock.Anything, "sync_launches", 0*time.Hour,
					mock.AnythingOfType("domain.SyncAction")).Run(func(args mock.Arguments) {
					fn := args.Get(3).(domain.SyncAction)
					fn(args.Get(0).(context.Context))
				}).Return(fmt.Errorf("failed to cancel bookings: something")).Once()
				spaceXQrsMock.On("GetUpcomingLaunches", mock.Anything).Return(upcomingLaunches, nil).Once()
				bookingsCmdsMock.On("Cancel", mock.Anything, mock.Anything).Return([]uuid.UUID{}, fmt.Errorf("something")).Once()
			},
		}, {
			name:    "failed to save launches",
			wantErr: fmt.Errorf("failed to save upcoming launches: something"),
			expectedMocks: func() {
				syncCmdsMock.On("SyncIfNecessary",
					mock.Anything, "sync_launches", 0*time.Hour,
					mock.AnythingOfType("domain.SyncAction")).Run(func(args mock.Arguments) {
					fn := args.Get(3).(domain.SyncAction)
					fn(args.Get(0).(context.Context))
				}).Return(fmt.Errorf("failed to save upcoming launches: something")).Once()
				spaceXQrsMock.On("GetUpcomingLaunches", mock.Anything).Return([]*domain.Launch{
					{LaunchPadID: "test", DateUTC: tnow.UTC()},
					{LaunchPadID: "test2", DateUTC: tnow.UTC().Add(24 * time.Hour)},
				}, nil).Once()
				bookingsCmdsMock.On("Cancel", mock.Anything, mock.Anything).Return([]uuid.UUID{}, nil).Once()
				mockedLogger.On("Debug", "%s bookings will be cancelled", 0, []uuid.UUID{})
				spaceXCmdsMock.On("SaveExternalLaunches", mock.Anything, upcomingLaunches).Return(fmt.Errorf("something")).Once()
			},
		}, {
			name:    "success",
			wantErr: nil,
			expectedMocks: func() {
				syncCmdsMock.On("SyncIfNecessary",
					mock.Anything, "sync_launches", 0*time.Hour,
					mock.AnythingOfType("domain.SyncAction")).Run(func(args mock.Arguments) {
					fn := args.Get(3).(domain.SyncAction)
					fn(args.Get(0).(context.Context))
				}).Return(nil).Once()
				spaceXQrsMock.On("GetUpcomingLaunches", mock.Anything).Return([]*domain.Launch{
					{LaunchPadID: "test", DateUTC: tnow.UTC()},
					{LaunchPadID: "test2", DateUTC: tnow.UTC().Add(24 * time.Hour)},
				}, nil).Once()
				bookingsCmdsMock.On("Cancel", mock.Anything, mock.Anything).Return([]uuid.UUID{}, nil).Once()
				mockedLogger.On("Debug", "%s bookings will be cancelled", 0, []uuid.UUID{})
				spaceXCmdsMock.On("SaveExternalLaunches", mock.Anything, upcomingLaunches).Return(nil).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewBookingsOrchestrator(mockedLogger, spaceXCmdsMock, spaceXQrsMock, bookingsCmdsMock, syncCmdsMock)

			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			err := h.SyncOnce(ctx, 0)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
		})
	}
}

func Test_handler_StartScheduledSync(t *testing.T) {
	mockedLogger := loggerMocks.NewInterface(t)
	bookingsCmdsMock := serviceMocks.NewBookingsServiceCommands(t)
	syncCmdsMock := serviceMocks.NewSyncServiceCommands(t)
	spaceXCmdsMock := serviceMocks.NewSpaceXServiceCommands(t)
	spaceXQrsMock := serviceMocks.NewSpaceXServiceQueries(t)

	ctx := context.Background()
	h := &handler{
		l:                mockedLogger,
		spaceXCommands:   spaceXCmdsMock,
		spaceXQueries:    spaceXQrsMock,
		bookingsCommands: bookingsCmdsMock,
		appsync:          syncCmdsMock,
		stopChan:         make(chan struct{}),
	}

	mockedLogger.On("Info", mock.Anything)
	syncCmdsMock.On("SyncIfNecessary",
		mock.Anything, "sync_launches", 100*time.Millisecond,
		mock.AnythingOfType("domain.SyncAction")).Return(nil)
	go h.StartScheduledSync(ctx, 100*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	// Stop the scheduler by sending to stopChan
	h.GracefulStop()

	mockedLogger.AssertExpectations(t)
	syncCmdsMock.AssertExpectations(t)

	h.stopChan = make(chan struct{})
	ctx1, cancel := context.WithCancel(ctx)

	go h.StartScheduledSync(ctx1, 100*time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	cancel()

	mockedLogger.AssertExpectations(t)
	syncCmdsMock.AssertExpectations(t)
}
