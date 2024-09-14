package launches

import (
	"context"
	"database/sql"
	"fmt"
	"spacet/internal/domain"
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"
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

// SaveLaunch creates a new launch in the database.
// Returns the created Launch's ID and an error in case of failure.
// If the launch already exists or a conflict is found, it returns domain.ErrLaunchAlreadyExists.
// If an internal error occurs, it logs the error and returns domain.ErrInternal.
func (r *CommandsRepo) CreateLaunch(ctx context.Context, launch domain.Launch) (id string, err error) {
	query := `
	INSERT INTO launches (external_id, domain, launch_name, date_utc, launchpad_id, destination, booking_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`
	err = r.db(ctx).QueryRow(ctx, query,
		launch.ID,
		launch.Domain,
		launch.Name,
		launch.DateUTC,
		launch.LaunchPadID,
		launch.Destination,
		sql.NullString{Valid: len(launch.BookingID) > 0, String: launch.BookingID.String()}).Scan(&id)
	if err != nil {
		if postgresql.IsConflictErr(err) {
			r.L.Debug(fmt.Errorf("launch %s already exists: %w", launch.ID, err))
			return id, domain.ErrLaunchAlreadyExists
		}
		r.L.Error(fmt.Errorf("failed to save launch: %w", err))
		return id, domain.ErrInternal
	}
	return id, nil
}

// SaveExternalLaunches creates a batch of launches in the database without associated booking.
// If an internal error occurs, it logs the error and returns domain.ErrInternal.
func (r *CommandsRepo) SaveExternalLaunches(ctx context.Context, launches []*domain.Launch) (err error) {
	// Prepare the insert query
	query := `INSERT INTO 
	launches (external_id, domain, launch_name, date_utc, launchpad_id, destination, lstatus)
	 VALUES %s ON CONFLICT (external_id) DO NOTHING`

	const nrArgs int = 7
	args := make([]interface{}, 0, len(launches)*nrArgs)
	where := ""
	for i, launch := range launches {
		r.L.Debug("upcoming %s on %s. ID: %s", launch.LaunchPadID, launch.DateUTC.Format("2006-01-02"), launch.ID)

		if i > 0 {
			where += ", "
		}
		where += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*nrArgs+1, i*nrArgs+2, i*nrArgs+3, i*nrArgs+4, i*nrArgs+5, i*nrArgs+6, i*nrArgs+7,
		)

		// Append the launch data to args
		args = append(args,
			launch.ExternalID,
			launch.Domain,
			launch.Name,
			launch.DateUTC,
			launch.LaunchPadID,
			launch.Destination,
			"scheduled",
		)
	}

	stmt := fmt.Sprintf(query, where)

	// Execute the batch insert
	_, err = r.db(ctx).Exec(ctx, stmt, args...)
	if err != nil {
		r.L.Error(fmt.Errorf("failed to execute batch insert: %w", err))
		return domain.ErrInternal
	}

	return nil
}
