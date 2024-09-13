package bookings

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"

	"github.com/google/uuid"
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

func (r *CommandsRepo) Cancel(ctx context.Context, restrictions []domain.LaunchRestriction) (cancelledBookings []uuid.UUID, _ error) {
	query := `UPDATE launches
		SET lstatus = 'cancelled'
		WHERE (trunc_date(date_utc), launchpad_id) IN (%s) AND booking_id IS NOT NULL
		RETURNING booking_id`

	// Build the query for batch update based on restrictions
	args := make([]interface{}, 0)
	where := ""
	for i, restriction := range restrictions {
		if i > 0 {
			where += ", "
		}
		where += fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2)
		args = append(args, restriction.DateUTC, restriction.LaunchPadID)
	}

	stmt := fmt.Sprintf(query, where)
	rows, err := r.db(ctx).Query(ctx, stmt, args...)
	if err != nil {
		r.L.Error(fmt.Errorf("failed to execute batch cancel query: %w", err))
		return nil, domain.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var bookingID uuid.UUID
		if err := rows.Scan(&bookingID); err != nil {
			r.L.Error(fmt.Errorf("failed to scan booking ID: %w", err))
			return nil, domain.ErrInternal
		}

		cancelledBookings = append(cancelledBookings, bookingID)
	}

	return cancelledBookings, nil
}
