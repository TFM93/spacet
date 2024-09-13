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
