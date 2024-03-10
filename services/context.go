package services

import (
	"context"

	"gorm.io/gorm"
)

type ContextKey string;

const (
	ContextKeyDB ContextKey = "gorm-db"
)

// The ContextKeyDB key must have a value associated with it inside the ctx.
func DB(ctx context.Context) *gorm.DB {
	//nolint:revive // See above.
	return ctx.Value(ContextKeyDB).(*gorm.DB) 
}
