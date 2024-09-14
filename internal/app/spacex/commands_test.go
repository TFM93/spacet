package spacex

import (
	"context"
	"fmt"
	domainMocks "spacet/gen/mocks/spacet/domain"
	loggerMocks "spacet/gen/mocks/spacet/pkg/logger"
	"spacet/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_commandsHandler_SaveExternalLaunches(t *testing.T) {
	tnow := time.Now()
	ctx := context.Background()
	mockedLogger := loggerMocks.NewInterface(t)
	lpadCommandsMock := domainMocks.NewLaunchPadRepoCommands(t)
	launchesCommandsMock := domainMocks.NewLaunchRepoCommands(t)
	spacexQueriesMock := domainMocks.NewSpaceXAPIQueries(t)
	mockedLogger.On("Debug", mock.Anything, mock.Anything)
	baseUUID := "2f506028-4a0e-47c4-826e-5a94db8a15e7"
	mars := domain.DestinationMars

	batchLauches := []*domain.Launch{}
	for i := 0; i <= 101; i++ {
		batchLauches = append(batchLauches, &domain.Launch{
			Name:        "something",
			Destination: &mars,
			DateUTC:     tnow.UTC(),
			BookingID:   uuid.MustParse(baseUUID),
		})
	}

	tests := []struct {
		name          string
		args          []*domain.Launch
		expectedMocks func()
		wantErr       error
	}{
		{
			name: "returns error",
			args: batchLauches,
			expectedMocks: func() {
				launchesCommandsMock.On("SaveExternalLaunches", mock.Anything, mock.Anything).Return(fmt.Errorf("something")).Once()
			},
			wantErr: fmt.Errorf("failed to save launches batch: something"),
		}, {
			name: "should batch",
			args: batchLauches,
			expectedMocks: func() {
				launchesCommandsMock.On("SaveExternalLaunches", mock.Anything, mock.Anything).Return(nil).Twice()
				mockedLogger.On("Info", "Launches update completed", "total", 102).Once()
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCommands(mockedLogger, spacexQueriesMock, lpadCommandsMock, launchesCommandsMock)

			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}

			err := h.SaveExternalLaunches(ctx, tt.args)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
		})
	}
}

func Test_commandsHandler_UpdateLaunchPads(t *testing.T) {
	ctx := context.Background()
	mockedLogger := loggerMocks.NewInterface(t)
	lpadCommandsMock := domainMocks.NewLaunchPadRepoCommands(t)
	launchesCommandsMock := domainMocks.NewLaunchRepoCommands(t)
	spacexQueriesMock := domainMocks.NewSpaceXAPIQueries(t)
	mockedLogger.On("Debug", mock.Anything)

	lpads := []*domain.LaunchPad{
		{
			Name:     "test",
			Locality: "asd",
			Region:   "ds",
			Timezone: "europe/lisbon",
			Status:   "active",
		}, {
			Name:     "test1",
			Locality: "asd1",
			Region:   "ds1",
			Timezone: "europe/lisbon",
			Status:   "retired",
		},
	}

	tests := []struct {
		name          string
		expectedMocks func()
		wantErr       error
	}{
		{
			name: "success",
			expectedMocks: func() {
				spacexQueriesMock.On("GetLaunchPads", mock.Anything).Return(lpads, nil).Once()
				lpadCommandsMock.On("SaveLaunchPad", mock.Anything, lpads[0]).Return("someid", nil).Once()
				lpadCommandsMock.On("SaveLaunchPad", mock.Anything, lpads[1]).Return("someid", nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "lpad persist error exits early",
			expectedMocks: func() {
				spacexQueriesMock.On("GetLaunchPads", mock.Anything).Return(lpads, nil).Once()
				lpadCommandsMock.On("SaveLaunchPad", mock.Anything, lpads[0]).Return("someid", fmt.Errorf("something")).Once()
			},
			wantErr: fmt.Errorf("something"),
		},
		{
			name: "lpad persist error exits early",
			expectedMocks: func() {
				spacexQueriesMock.On("GetLaunchPads", mock.Anything).Return(lpads, fmt.Errorf("something")).Once()
			},
			wantErr: fmt.Errorf("something"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCommands(mockedLogger, spacexQueriesMock, lpadCommandsMock, launchesCommandsMock)

			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			err := h.UpdateLaunchPads(ctx)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
		})
	}
}
