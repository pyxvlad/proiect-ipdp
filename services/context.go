package services

import (
	"context"

	"github.com/pyxvlad/proiect-ipdp/database"
)

type ContextKey string

const (
	ContextKeyDB ContextKey = "db-txlike"
)

// The ContextKeyDB key must have a value associated with it inside the ctx.
func DBTxLike(ctx context.Context) database.TxLike {
	//nolint:revive // See above.
	return ctx.Value(ContextKeyDB).(database.TxLike)
}

func DB(ctx context.Context) *database.Queries {
	return database.New(DBTxLike(ctx))
}
