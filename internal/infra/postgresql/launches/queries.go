package launches

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"
	"time"
)

type QueriesRepo struct {
	PG postgresql.Interface
	L  log.Interface
}

func (r QueriesRepo) db(ctx context.Context) postgresql.DBProvider {
	tx, ok := ctx.Value(domain.TxKey).(postgresql.Tx)
	if ok {
		return tx
	}
	return r.PG.GetPool()
}

// IsLaunchpadAvailableForDate checks if a given launchpad is available on the specified date.
// It expects that the implementation locks the launches for updating
func (r *QueriesRepo) IsLaunchpadAvailableForDate(ctx context.Context, launchpadID string, date time.Time) (bool, error) {
	rows, err := r.db(ctx).Query(ctx, `
        SELECT id
        FROM launches
        WHERE launchpad_id = $1
        AND date_utc::date = $2::date
        AND lstatus != 'cancelled'
        FOR UPDATE
    `, launchpadID, date)

	if err != nil {
		r.L.Error(fmt.Errorf("failed to check launchpad availability: %w", err))
		return false, domain.ErrInternal
	}
	defer rows.Close()

	return !rows.Next(), nil
}

// LaunchesOnSameDestinationOnTargetWeek counts the number of launches for the date's week, launchpad and destination.
// It expects that the implementation locks the launches for updating
// If an internal error occurs, it logs the error and returns domain.ErrInternal.
func (r *QueriesRepo) LaunchesOnSameDestinationOnTargetWeek(ctx context.Context, launchpadID string, date time.Time, destination string) (count int, err error) {
	query := `
        SELECT id
        FROM launches
        WHERE launchpad_id = $1
        AND date_utc >= date_trunc('week'::text, $2::timestamp)
        AND date_utc < date_trunc('week'::text, $2::timestamp) + interval '7 days'
        AND destination = $3
        AND lstatus != 'cancelled'
        FOR UPDATE
    `

	rows, err := r.db(ctx).Query(ctx, query, launchpadID, date, destination)
	if err != nil {
		r.L.Error(fmt.Errorf("failed to get launches destinations for the week, for update: %w", err))
		return 0, domain.ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		count++
	}

	return count, nil
}
