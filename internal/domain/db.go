package domain

import (
	"context"
)

// Transaction is an interface for managing database transactions
type Transaction interface {
	BeginTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type ctxKey int

const (
	// TxKey is used to save the transaction into the context
	// iota type is the recommendation here in order to avoid key collisions
	TxKey ctxKey = iota
)
