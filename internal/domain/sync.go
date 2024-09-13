package domain

import (
	"context"
	"time"
)

// SyncRepoCommands is an interface for distributed locking and sync info updates
type SyncRepoCommands interface {
	TryDistributedLock(ctx context.Context, key uint32) (bool, error)
	ReleaseDistributedLock(ctx context.Context, key uint32) error
	GetLastSyncTimestamp(ctx context.Context, resourceName string) (time.Time, error)
	UpdateLastSyncTimestamp(ctx context.Context, resourceName string, newTimestamp time.Time) error
}

type LockKeyGenerator func(resourceName string) uint32

type SyncAction func(ctx context.Context) error
