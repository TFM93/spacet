package healthcheck

import (
	"context"
	dbmocks "spacet/gen/mocks/spacet/pkg/postgresql"
	"spacet/internal/domain"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_healthCheckQueries_Check(t *testing.T) {
	type fields struct {
		repo domain.MonitoringInfraQueries
	}
	type args struct {
		ctx context.Context
	}
	mockDB := dbmocks.NewInterface(t)

	tests := []struct {
		name          string
		fields        fields
		expectedMocks func()
		args          args
		want          bool
	}{
		{
			name: "success",
			fields: fields{
				repo: mockDB,
			},
			expectedMocks: func() {
				mockDB.On("Ping", mock.Anything).Return(true).Once()
			},
			args: args{
				ctx: context.Background(),
			},
			want: true,
		},
		{
			name: "unhealthy repo",
			fields: fields{
				repo: mockDB,
			},
			expectedMocks: func() {
				mockDB.On("Ping", mock.Anything).Return(false).Once()
			},
			args: args{
				ctx: context.Background(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			h := NewQueries(tt.fields.repo)
			if got := h.Check(tt.args.ctx); got != tt.want {
				t.Errorf("healthCheckQueries.Check() = %v, want %v", got, tt.want)
			}

		})
	}
}
