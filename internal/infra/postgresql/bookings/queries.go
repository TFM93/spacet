package bookings

import (
	"context"
	"fmt"
	"spacet/internal/domain"
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"
	"strings"
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

// ListTickets fetches the booked launches from the database based on the provided cursor data
// If the query fails to execute, it returns return domain.ErrInternal
// If theres an error processing the data, it returns domain.ErrFailedToProcessData
func (r QueriesRepo) ListTickets(ctx context.Context, cursorTicketID string, cursorUpdatedAt *time.Time, limit int32, filters domain.TicketSearchFilters) ([]*domain.Ticket, error) {
	var whereClauses []string
	var args []any

	// pagination
	if cursorTicketID != "" && cursorUpdatedAt != nil {
		whereClauses = append(whereClauses, `(l.updated_at < $1 OR (l.updated_at = $1 AND l.booking_id < $2))`)
		args = append(args, *cursorUpdatedAt, cursorTicketID)
	}

	if filters.FirstName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("b.first_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*filters.FirstName+"%")
	}
	if filters.LastName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("b.last_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*filters.LastName+"%")
	}
	if filters.Destination != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("l.destination ILIKE $%d", len(args)+1))
		args = append(args, "%"+*filters.Destination+"%")
	}
	if filters.LaunchPadID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("l.launchpad_id ILIKE $%d", len(args)+1))
		args = append(args, "%"+*filters.LaunchPadID+"%")
	}
	if filters.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("l.lstatus ILIKE $%d", len(args)+1))
		args = append(args, "%"+*filters.Status+"%")
	}

	where := ""
	if len(whereClauses) > 0 {
		where = "WHERE " + strings.Join(whereClauses, " AND ")
	}
	query := fmt.Sprintf(`
	select b.id, b.first_name, b.last_name, l.launchpad_id, l.date_utc, l.lstatus, l.destination, l.created_at, l.updated_at
	FROM launches l INNER JOIN bookings b ON l.booking_id = b.id 
	%s 
	ORDER BY l.updated_at DESC, l.id DESC LIMIT $%d`, where, len(args)+1)
	args = append(args, limit)

	rows, err := r.db(ctx).Query(ctx, query, args...)
	if err != nil {
		r.L.Debug(fmt.Errorf("failed to list booked launches: %w", err))
		return nil, domain.ErrInternal
	}
	defer rows.Close()
	var tickets []*domain.Ticket
	// NOTE: as per version 5 of pgx this can be done relying on generics:
	// https://donchev.is/post/working-with-postgresql-in-go-using-pgx/
	for rows.Next() {
		var ticket domain.Ticket
		if err := rows.Scan(
			&ticket.ID,
			&ticket.FirstName,
			&ticket.LastName,
			&ticket.LaunchPadID,
			&ticket.LaunchDate,
			&ticket.Status,
			&ticket.Destination,
			&ticket.CreatedAt,
			&ticket.UpdatedAt); err != nil {
			r.L.Error(fmt.Errorf("failed to scan row: %w", err))
			return nil, domain.ErrFailedToProcessData
		}
		tickets = append(tickets, &ticket)
	}

	if err := rows.Err(); err != nil {
		r.L.Error(fmt.Errorf("row iteration error: %w", err))
		return nil, domain.ErrFailedToProcessData
	}
	return tickets, nil

}
