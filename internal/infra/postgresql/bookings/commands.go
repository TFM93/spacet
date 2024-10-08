package bookings

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"
	"strings"
	"time"

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

func (r *CommandsRepo) CancelByID(ctx context.Context, bookingID uuid.UUID) error {
	query := `
			UPDATE launches
			SET lstatus = 'cancelled'
			WHERE booking_id = $1
		`
	commandTag, err := r.db(ctx).Exec(ctx, query, bookingID)
	if err != nil {
		r.L.Error(fmt.Errorf("failed to cancel booking: %w", err))
		return domain.ErrInternal
	}
	if commandTag.RowsAffected() == 0 {
		r.L.Debug("launch with bookingID %s does not exist", bookingID)
		return domain.ErrBookingNotFound
	}
	return nil
}

// Cancel changes the status of the launches to cancelled if they meet the launchPad daily restriction (only applies to the launches with booking)
func (r *CommandsRepo) Cancel(ctx context.Context, restrictions map[string][]time.Time) (cancelledBookings []uuid.UUID, _ error) {
	for launchPad, dates := range restrictions {
		cancelled, err := r.cancelPerLaunchPad(ctx, launchPad, dates)
		if err != nil {
			r.L.Error(fmt.Errorf("failed to cancel per launchpad query: %w", err))
			return nil, domain.ErrInternal
		}
		cancelledBookings = append(cancelledBookings, cancelled...)
	}
	return
}

func (r *CommandsRepo) cancelPerLaunchPad(ctx context.Context, launchPadID string, dates []time.Time) (cancelledBookings []uuid.UUID, _ error) {
	query := `
			UPDATE launches
			SET lstatus = 'cancelled'
			WHERE id IN (
				SELECT id
				FROM launches
				WHERE domain = 'SPACET'
				AND lstatus != 'cancelled'
				AND launchpad_id = $%d
				AND trunc_date(date_utc) IN (%s)
				FOR UPDATE
			)
			RETURNING booking_id
		`

	var conditions []string
	var args []interface{}
	argIndex := 1

	for _, date := range dates {
		conditions = append(conditions, fmt.Sprintf("trunc_date($%d)", argIndex))
		args = append(args, date)
		argIndex++
	}
	args = append(args, launchPadID)
	whereClause := strings.Join(conditions, ",")
	stmt := fmt.Sprintf(query, argIndex, whereClause)

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
	return
}

// CreateBooking persists the provided booking and returns the autogenerated booking id
func (r *CommandsRepo) CreateBooking(ctx context.Context, booking domain.Booking) (id uuid.UUID, err error) {
	query := `
	INSERT INTO bookings (first_name, last_name, gender, birthday)
	VALUES ($1, $2, $3, $4)
	RETURNING id
`
	err = r.db(ctx).QueryRow(ctx, query,
		booking.FirstName,
		booking.LastName,
		booking.Gender,
		booking.BirthDay).Scan(&id)
	if err != nil {
		r.L.Error(fmt.Errorf("failed to save booking: %w", err))
		return id, domain.ErrInternal
	}
	return id, nil
}
