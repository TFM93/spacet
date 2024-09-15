package bookings

import (
	"context"
	"fmt"
	"reflect"
	loggermocks "spacet/gen/mocks/spacet/pkg/logger"
	dbmocks "spacet/gen/mocks/spacet/pkg/postgresql"
	"spacet/internal/domain"
	"spacet/pkg/postgresql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func TestCommandsRepo_Cancel(t *testing.T) {
	ctx := context.Background()
	mockLogger := loggermocks.NewInterface(t)

	mockDB := dbmocks.NewInterface(t)

	tests := []struct {
		name                  string
		restrictions          map[string][]time.Time
		expectedMocks         func()
		wantCancelledBookings []uuid.UUID
		wantErr               error
	}{
		{
			name:                  "without restrictions - noop",
			restrictions:          make(map[string][]time.Time),
			wantCancelledBookings: []uuid.UUID(nil),
			wantErr:               nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			r := &CommandsRepo{
				PG: mockDB,
				L:  mockLogger,
			}
			gotCancelledBookings, err := r.Cancel(ctx, tt.restrictions)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if !reflect.DeepEqual(gotCancelledBookings, tt.wantCancelledBookings) {
				t.Errorf("CommandsRepo.Cancel() = %v, want %v", gotCancelledBookings, tt.wantCancelledBookings)
			}
		})
	}
}

func TestCommandsRepo_CreateBooking(t *testing.T) {
	expectedBookingID := "0f913f6a-497b-4305-b3d1-3f53657e3a25"
	ctx := context.Background()
	tnow := time.Now()
	mockDB := dbmocks.NewInterface(t)
	mockLogger := loggermocks.NewInterface(t)
	mockDBProvider := dbmocks.NewDBProvider(t)
	mockRow := new(MockRow)
	r := CommandsRepo{
		PG: mockDB,
		L:  mockLogger,
	}
	mockDB.On("GetPool").Return(mockDBProvider)

	tests := []struct {
		name          string
		booking       domain.Booking
		expectedMocks func()
		wantId        uuid.UUID
		wantErr       error
	}{
		{
			name: "create booking",
			booking: domain.Booking{
				FirstName: "first",
				LastName:  "last",
				BirthDay:  tnow,
				Gender:    domain.GenderFemale,
			},
			expectedMocks: func() {
				mockDBProvider.On("QueryRow", mock.Anything,
					"\n\tINSERT INTO bookings (first_name, last_name, gender, birthday)\n\tVALUES ($1, $2, $3, $4)\n\tRETURNING id\n",
					"first", "last", domain.GenderFemale, tnow).Return(mockRow).Once()
				mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
					// change the value of the scan argument
					arg := args.Get(0).([]interface{})
					*arg[0].(*uuid.UUID) = uuid.MustParse(expectedBookingID)
				}).Return(nil).Once()
			},
			wantId:  uuid.MustParse(expectedBookingID),
			wantErr: nil,
		}, {
			name: "create booking- failure",
			booking: domain.Booking{
				FirstName: "first",
				LastName:  "last",
				BirthDay:  tnow,
				Gender:    domain.GenderFemale,
			},
			expectedMocks: func() {
				mockDBProvider.On("QueryRow", mock.Anything,
					"\n\tINSERT INTO bookings (first_name, last_name, gender, birthday)\n\tVALUES ($1, $2, $3, $4)\n\tRETURNING id\n",
					"first", "last", domain.GenderFemale, tnow).Return(mockRow).Once()
				mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
					// change the value of the scan argument
					arg := args.Get(0).([]interface{})
					*arg[0].(*uuid.UUID) = uuid.MustParse(expectedBookingID)
				}).Return(fmt.Errorf("something")).Once()
				mockLogger.On("Error", mock.Anything).Return()

			},
			wantId:  uuid.MustParse(expectedBookingID),
			wantErr: domain.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			gotId, err := r.CreateBooking(ctx, tt.booking)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if gotId != tt.wantId {
				t.Errorf("CreateBooking() = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestCommandsRepo_db(t *testing.T) {

	mockDB := dbmocks.NewInterface(t)
	mockLogger := loggermocks.NewInterface(t)
	mockTx := dbmocks.NewTx(t)
	mockDBProvider := dbmocks.NewDBProvider(t)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name          string
		args          args
		expectedMocks func()
		want          postgresql.DBProvider
	}{
		{
			name: "return tx",
			args: args{
				ctx: context.WithValue(context.Background(), domain.TxKey, mockTx),
			},
			want: mockTx,
		},
		{
			name: "return from pool",
			args: args{
				ctx: context.Background(),
			},
			expectedMocks: func() {
				mockDB.On("GetPool").Return(mockDBProvider).Once()
			},
			want: mockDBProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CommandsRepo{
				PG: mockDB,
				L:  mockLogger,
			}
			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			if got := r.db(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandsRepo.db() = %v, want %v", got, tt.want)
			}
		})
	}
}
