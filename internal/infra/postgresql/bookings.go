package postgresql

import (
	"spacet/internal/domain"
	"spacet/internal/infra/postgresql/bookings"
	"spacet/pkg/logger"
	"spacet/pkg/postgresql"
)

// NewBookingQueriesRepo creates a new instance of bookingQueriesRepo that satisfies the domain.BookingRepoQueries interface
func NewBookingQueriesRepo(pg postgresql.Interface, logger logger.Interface) domain.BookingRepoQueries {
	ur := &bookings.QueriesRepo{PG: pg, L: logger}
	return ur
}

// NewBookingCommandsRepo creates a new instance of commandsRepo that satisfies the domain.BookingRepoCommands interface
func NewBookingCommandsRepo(pg postgresql.Interface, logger logger.Interface) domain.BookingRepoCommands {
	ur := &bookings.CommandsRepo{PG: pg, L: logger}
	return ur
}
