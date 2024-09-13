package postgresql

import (
	"spacet/internal/domain"
	"spacet/internal/infra/postgresql/launchpads"
	"spacet/pkg/logger"
	"spacet/pkg/postgresql"
)

// NewLaunchPadCommandsRepo creates a new instance of commandsRepo that satisfies the domain.LaunchPadRepoCommands interface
func NewLaunchPadCommandsRepo(pg postgresql.Interface, logger logger.Interface) domain.LaunchPadRepoCommands {
	ur := &launchpads.CommandsRepo{PG: pg, L: logger}
	return ur
}
