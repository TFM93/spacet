package bookings

import (
	"context"
	"fmt"
	"spacet/pkg/logger"
)

type Commands interface {
	// TODO: describe and implement
	Cancel(ctx context.Context) (err error)
}

type handler struct {
	l logger.Interface
}

func NewCommands(logger logger.Interface) Commands {
	return &handler{l: logger}
}

func (h handler) Cancel(ctx context.Context) (err error) {
	return fmt.Errorf("not implemented yet")
}
