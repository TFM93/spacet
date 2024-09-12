package bookings

import (
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"
)

type CommandsRepo struct {
	PG postgresql.Interface
	L  log.Interface
}
