package sync

import (
	"context"
	"spacet/internal/domain"
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"
	"time"
)

type CommandsRepo struct {
	PG postgresql.Interface
	L  log.Interface
}

func (r CommandsRepo) db(ctx context.Context) postgresql.DBProvider {
	tx, ok := ctx.Value(domain.TxKey).(postgresql.Tx)
	if ok {
		return tx
	}
	return r.PG.GetPool()
}

func (r *CommandsRepo) TryDistributedLock(ctx context.Context, key uint32) (bool, error) {
	var locked bool
	err := r.db(ctx).QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", key).Scan(&locked)
	return locked, err
}

func (r *CommandsRepo) ReleaseDistributedLock(ctx context.Context, key uint32) error {
	_, err := r.db(ctx).Exec(ctx, "SELECT pg_advisory_unlock($1)", key)
	return err
}

func (r *CommandsRepo) GetLastSyncTimestamp(ctx context.Context, resourceName string) (time.Time, error) {
	var lastSync time.Time
	err := r.db(ctx).QueryRow(ctx, "SELECT last_sync FROM sync_info WHERE resource_name = $1", resourceName).Scan(&lastSync)
	if err == postgresql.ErrNoRows {
		return time.Time{}, nil
	}
	return lastSync, err
}

func (r *CommandsRepo) UpdateLastSyncTimestamp(ctx context.Context, resourceName string, newTimestamp time.Time) error {
	commandTag, err := r.db(ctx).Exec(ctx, "UPDATE sync_info SET last_sync = $1 WHERE resource_name = $2", newTimestamp, resourceName)
	if commandTag.RowsAffected() == 0 {
		_, err = r.db(ctx).Exec(ctx, "INSERT INTO sync_info(resource_name, last_sync) VALUES ($1, $2)", resourceName, newTimestamp)
	}
	return err
}
