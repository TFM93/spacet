package postgresql

import (
	"spacet/internal/domain"
	"spacet/internal/infra/postgresql/sync"
	"spacet/pkg/logger"
	"spacet/pkg/postgresql"
)

// NewSyncCommandsRepo creates a new instance of commandsRepo that satisfies the domain.SyncRepoCommands interface
func NewSyncCommandsRepo(pg postgresql.Interface, logger logger.Interface) domain.SyncRepoCommands {
	ur := &sync.CommandsRepo{PG: pg, L: logger}
	return ur
}
