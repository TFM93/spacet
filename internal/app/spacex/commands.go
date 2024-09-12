package spacex

import (
	"context"
	"fmt"
	"spacet/pkg/logger"
	"time"
)

type Commands interface {
	// Handle fetches data from spaceX api and persists the information. It also checks for conflicts and cancels bookings if necessary
	// TODO: describe error handling
	Handle(ctx context.Context, req SyncSpaceXDataCommand) (err error)
}

type SyncSpaceXDataCommand struct {
	SyncTime time.Time
}

type dataHandler struct {
	l logger.Interface
}

func NewCommands(logger logger.Interface) Commands {
	return &dataHandler{l: logger}
}

func (h dataHandler) Handle(ctx context.Context, req SyncSpaceXDataCommand) (err error) {
	return fmt.Errorf("not implemented yet")
}
