package sync

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	"spacet/pkg/logger"
	"time"
)

type Commands interface {
	// SyncIfNecessary checks if a month has passed and updates the resource if needed
	SyncIfNecessary(ctx context.Context, resourceName string, syncInterval time.Duration, syncFn domain.SyncAction) error
}

type handler struct {
	l                logger.Interface
	transaction      domain.Transaction
	syncRepo         domain.SyncRepoCommands
	lockKeyGenerator domain.LockKeyGenerator
}

func NewCommands(logger logger.Interface, transaction domain.Transaction, repo domain.SyncRepoCommands, lockKeyGenerator domain.LockKeyGenerator) Commands {
	return &handler{l: logger, transaction: transaction, syncRepo: repo, lockKeyGenerator: lockKeyGenerator}
}

// SyncIfNecessary checks if a month has passed and updates the resource if needed
func (h *handler) SyncIfNecessary(ctx context.Context, resourceName string, syncInterval time.Duration, syncFn domain.SyncAction) error {
	lockKey := h.lockKeyGenerator(resourceName)
	locked, err := h.syncRepo.TryDistributedLock(ctx, lockKey)
	if err != nil {
		return fmt.Errorf("failed to acquire advisory lock: %w", err)
	}
	if !locked {
		// Another instance holds the lock, skip the sync
		return nil
	}
	defer h.syncRepo.ReleaseDistributedLock(ctx, lockKey)

	return h.transaction.BeginTx(ctx, func(ctx context.Context) error {
		lastSync, err := h.syncRepo.GetLastSyncTimestamp(ctx, resourceName)
		if err != nil {
			return fmt.Errorf("failed to fetch last sync timestamp: %w", err)
		}

		utcNow := time.Now().UTC()
		if utcNow.Sub(lastSync) >= syncInterval {
			err = syncFn(ctx)
			if err != nil {
				return fmt.Errorf("sync action failed: %w", err)
			}

			err = h.syncRepo.UpdateLastSyncTimestamp(ctx, resourceName, utcNow)
			if err != nil {
				return fmt.Errorf("failed to update last sync timestamp: %w", err)
			}
		}

		return nil
	})

}
