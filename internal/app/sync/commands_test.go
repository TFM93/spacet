package sync

import (
	"context"
	"fmt"
	domainMocks "spacet/gen/mocks/spacet/domain"
	loggerMocks "spacet/gen/mocks/spacet/pkg/logger"
	"spacet/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_handler_SyncIfNecessary(t *testing.T) {
	tnow := time.Now()
	ctx := context.Background()
	mockedLogger := loggerMocks.NewInterface(t)
	transactionMock := domainMocks.NewTransaction(t)
	syncRepoMock := domainMocks.NewSyncRepoCommands(t)
	lockKeyGenMock := domainMocks.NewLockKeyGenerator(t)

	type args struct {
		resourceName string
		syncInterval time.Duration
		syncFn       domain.SyncAction
	}
	tests := []struct {
		name          string
		args          args
		expectedMocks func()
		wantErr       error
	}{
		{
			name: "error getting lock",
			args: args{
				resourceName: "test",
				syncInterval: time.Hour,
				syncFn: func(ctx context.Context) error {
					return nil
				},
			},
			expectedMocks: func() {
				lockKeyGenMock.On("Execute", "test").Return(uint32(1)).Once()
				syncRepoMock.On("TryDistributedLock", mock.Anything, uint32(1)).Return(false, fmt.Errorf("something")).Once()
			},
			wantErr: fmt.Errorf("failed to acquire advisory lock: something"),
		}, {
			name: "another instance holds the lock",
			args: args{
				resourceName: "test",
				syncInterval: time.Hour,
				syncFn: func(ctx context.Context) error {
					return nil
				},
			},
			expectedMocks: func() {
				lockKeyGenMock.On("Execute", "test").Return(uint32(1)).Once()
				syncRepoMock.On("TryDistributedLock", mock.Anything, uint32(1)).Return(false, nil).Once()
			},
			wantErr: nil,
		}, {
			name: "failed to fetch last sync timestamp",
			args: args{
				resourceName: "test",
				syncInterval: time.Hour,
				syncFn: func(ctx context.Context) error {
					return nil
				},
			},
			expectedMocks: func() {
				lockKeyGenMock.On("Execute", "test").Return(uint32(1)).Once()
				syncRepoMock.On("TryDistributedLock", mock.Anything, uint32(1)).Return(true, nil).Once()
				syncRepoMock.On("ReleaseDistributedLock", mock.Anything, uint32(1)).Return(nil).Once()
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(fmt.Errorf("failed to fetch last sync timestamp: something")).Once()
				syncRepoMock.On("GetLastSyncTimestamp", mock.Anything, "test").Return(time.Time{}, fmt.Errorf("something")).Once()

			},
			wantErr: fmt.Errorf("failed to fetch last sync timestamp: something"),
		}, {
			name: "failed to execute syncFn",
			args: args{
				resourceName: "test",
				syncInterval: time.Hour,
				syncFn: func(ctx context.Context) error {
					return fmt.Errorf("something")
				},
			},
			expectedMocks: func() {
				lockKeyGenMock.On("Execute", "test").Return(uint32(1)).Once()
				syncRepoMock.On("TryDistributedLock", mock.Anything, uint32(1)).Return(true, nil).Once()
				syncRepoMock.On("ReleaseDistributedLock", mock.Anything, uint32(1)).Return(nil).Once()
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(fmt.Errorf("sync action failed: something")).Once()
				syncRepoMock.On("GetLastSyncTimestamp", mock.Anything, "test").Return(tnow.UTC().Add(-2*time.Hour), nil).Once()

			},
			wantErr: fmt.Errorf("sync action failed: something"),
		}, {
			name: "failed to update last sync timestamp",
			args: args{
				resourceName: "test",
				syncInterval: time.Hour,
				syncFn: func(ctx context.Context) error {
					return nil
				},
			},
			expectedMocks: func() {
				lockKeyGenMock.On("Execute", "test").Return(uint32(1)).Once()
				syncRepoMock.On("TryDistributedLock", mock.Anything, uint32(1)).Return(true, nil).Once()
				syncRepoMock.On("ReleaseDistributedLock", mock.Anything, uint32(1)).Return(nil).Once()
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(fmt.Errorf("failed to update last sync timestamp: something")).Once()
				syncRepoMock.On("GetLastSyncTimestamp", mock.Anything, "test").Return(tnow.UTC().Add(-2*time.Hour), nil).Once()
				syncRepoMock.On("UpdateLastSyncTimestamp", mock.Anything, "test", mock.Anything).Return(fmt.Errorf("something")).Once()

			},
			wantErr: fmt.Errorf("failed to update last sync timestamp: something"),
		}, {
			name: "success",
			args: args{
				resourceName: "test",
				syncInterval: time.Hour,
				syncFn: func(ctx context.Context) error {
					return nil
				},
			},
			expectedMocks: func() {
				lockKeyGenMock.On("Execute", "test").Return(uint32(1)).Once()
				syncRepoMock.On("TryDistributedLock", mock.Anything, uint32(1)).Return(true, nil).Once()
				syncRepoMock.On("ReleaseDistributedLock", mock.Anything, uint32(1)).Return(nil).Once()
				transactionMock.On("BeginTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					fn(args.Get(0).(context.Context))
				}).Return(nil).Once()
				syncRepoMock.On("GetLastSyncTimestamp", mock.Anything, "test").Return(tnow.UTC().Add(-2*time.Hour), nil).Once()
				syncRepoMock.On("UpdateLastSyncTimestamp", mock.Anything, "test", mock.Anything).Return(nil).Once()

			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCommands(mockedLogger, transactionMock, syncRepoMock, lockKeyGenMock.Execute)
			if tt.expectedMocks != nil {
				tt.expectedMocks()
			}
			err := h.SyncIfNecessary(ctx, tt.args.resourceName, tt.args.syncInterval, tt.args.syncFn)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
		})
	}
}
