package bookings

import (
	"context"
	"spacet/internal/domain"
	"spacet/pkg/logger"

	"github.com/google/uuid"
)

type Commands interface {
	// TODO: describe and implement
	Cancel(ctx context.Context, restriction []domain.LaunchRestriction) (cancelled []uuid.UUID, err error)
}

type handler struct {
	l    logger.Interface
	repo domain.BookingRepoCommands
}

func NewCommands(logger logger.Interface, repo domain.BookingRepoCommands) Commands {
	return &handler{l: logger, repo: repo}
}

func (h handler) Cancel(ctx context.Context, restriction []domain.LaunchRestriction) (cancelled []uuid.UUID, err error) {
	return h.repo.Cancel(ctx, restriction)
}
