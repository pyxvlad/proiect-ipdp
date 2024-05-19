package services

import (
	"context"

	"github.com/pyxvlad/proiect-ipdp/database"
)

type ContextKey string

const (
	ContextKeyDB        ContextKey = "db-txlike"
	ContextKeyCoverPath ContextKey = "cover-path"
)

// The ContextKeyDB key must have a value associated with it inside the ctx.
func DBTxLike(ctx context.Context) database.TxLike {
	//nolint:revive // See above.
	return ctx.Value(ContextKeyDB).(database.TxLike)
}

func DB(ctx context.Context) *database.Queries {
	return database.New(DBTxLike(ctx))
}

// The ContextKeyDB key must have a value associated with it inside the ctx.
func CoverPath(ctx context.Context) string {
	//nolint:revive // See above.
	return ctx.Value(ContextKeyCoverPath).(string)
}
