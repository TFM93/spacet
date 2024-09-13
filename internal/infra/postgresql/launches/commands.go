package launches

import (
	"context"
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

func (r *CommandsRepo) SaveLaunch(ctx context.Context, launch *domain.Launch) (id string, err error) {
	query := `
	INSERT INTO launches (external_id, domain, launch_name, date_utc, launchpad_id, destination)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	err = r.db(ctx).QueryRow(ctx, query,
		launch.ID,
		launch.Domain,
		launch.Name,
		launch.DateUTC,
		launch.LaunchPadID,
		launch.Destination).Scan(&id)
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
