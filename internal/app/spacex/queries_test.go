package spacex

import (
	"context"
	"reflect"
	domainMocks "spacet/gen/mocks/spacet/domain"
	loggerMocks "spacet/gen/mocks/spacet/pkg/logger"
	"spacet/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_queriesHandler_GetUpcomingLaunches(t *testing.T) {
	ctx := context.Background()
	mockedLogger := loggerMocks.NewInterface(t)
	spacexQueriesMock := domainMocks.NewSpaceXAPIQueries(t)
	mars := domain.DestinationMars

	baseLaunch := domain.Launch{
		ID:          "something",
		ExternalID:  "else",
		Domain:      domain.SpaceXDomain,
		Name:        "special",
		DateUTC:     time.Now().UTC(),
		LaunchPadID: "123",
		Destination: &mars,
		Status:      "scheduled",
		BookingID:   uuid.MustParse("2f506028-4a0e-47c4-826e-5a94db8a15e7"),
	}
	tests := []struct {
		name          string
		expectedMocks func()
		want          []*domain.Launch
		wantErr       error
	}{
		{
			name: "success",
			expectedMocks: func() {
				spacexQueriesMock.On("GetUpcomingLaunches", ctx).Return(
					[]*domain.Launch{&baseLaunch}, nil).Once()
			},
			want:    []*domain.Launch{&baseLaunch},
			wantErr: nil,
		},
		{
			name: "withErr",
			expectedMocks: func() {
				spacexQueriesMock.On("GetUpcomingLaunches", ctx).Return(
					[]*domain.Launch{}, domain.ErrInternal).Once()
			},
			want:    []*domain.Launch{},
			wantErr: domain.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewQueries(mockedLogger, spacexQueriesMock)
			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}

			got, err := h.GetUpcomingLaunches(ctx)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queriesHandler.GetUpcomingLaunches() = %v, want %v", got, tt.want)
			}
		})
	}
}
