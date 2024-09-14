package postgresql

import (
	"spacet/internal/domain"
	"spacet/internal/infra/postgresql/launches"
	"spacet/pkg/logger"
	"spacet/pkg/postgresql"
)

// NewLaunchesCommandsRepo creates a new instance of commandsRepo that satisfies the domain.LaunchRepoCommands interface
func NewLaunchesCommandsRepo(pg postgresql.Interface, logger logger.Interface) domain.LaunchRepoCommands {
	ur := &launches.CommandsRepo{PG: pg, L: logger}
	return ur
}

// NewLaunchesQueriesRepo creates a new instance of queriesRepo that satisfies the domain.LaunchRepoQueries interface
func NewLaunchesQueriesRepo(pg postgresql.Interface, logger logger.Interface) domain.LaunchRepoQueries {
	ur := &launches.QueriesRepo{PG: pg, L: logger}
	return ur
}
