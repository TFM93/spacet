package bookings

import (
	log "spacet/pkg/logger"
	"spacet/pkg/postgresql"
)

type QueriesRepo struct {
	PG postgresql.Interface
	L  log.Interface
}
