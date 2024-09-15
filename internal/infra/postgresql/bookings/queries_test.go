package bookings

import (
	"context"
	"reflect"
	loggermocks "spacet/gen/mocks/spacet/pkg/logger"
	dbmocks "spacet/gen/mocks/spacet/pkg/postgresql"
	"spacet/internal/domain"
	"spacet/pkg/postgresql"
	"testing"
)

func TestQueriesRepo_db(t *testing.T) {

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
			r := &QueriesRepo{
				PG: mockDB,
				L:  mockLogger,
			}
			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			if got := r.db(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueriesRepo.db() = %v, want %v", got, tt.want)
			}
		})
	}
}
