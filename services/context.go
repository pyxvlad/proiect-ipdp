package services

import (
	"context"

	"gorm.io/gorm"
)

type ContextKey string;

const (
	ContextKeyDB ContextKey = "gorm-db"
)

func DB(ctx context.Context) *gorm.DB {
	value := ctx.Value(ContextKeyDB)
	return value.(*gorm.DB)
}
