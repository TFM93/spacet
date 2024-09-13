package launchpads

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

// SaveLaunchPad creates a new launch in the database.
// Returns the created LaunchpPad's ID and an error in case of failure.
// If the launchpad already exists or a conflict is found, it tries to update the launchpad.
// If an internal error occurs, it logs the error and returns domain.ErrInternal.
func (r CommandsRepo) SaveLaunchPad(ctx context.Context, launchPad *domain.LaunchPad) (id string, err error) {
	query := `
		INSERT INTO launchpads (id, pad_name, locality, region, timezone, pad_status)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			pad_name = EXCLUDED.pad_name,
			locality = EXCLUDED.locality,
			region = EXCLUDED.region,
			timezone = EXCLUDED.timezone,
			pad_status = EXCLUDED.pad_status
		RETURNING id
	`
	err = r.db(ctx).QueryRow(ctx, query,
		launchPad.ID,
		launchPad.Name,
		launchPad.Locality,
		launchPad.Region,
		launchPad.Timezone,
		launchPad.Status).Scan(&id)
	if err != nil {
		r.L.Error(fmt.Errorf("failed to save launchpad: %w", err))
		return id, domain.ErrInternal
	}
	return id, nil
}
