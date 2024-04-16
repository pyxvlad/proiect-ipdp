package database

//go:generate sqlc generate

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	gooseDB "github.com/pressly/goose/v3/database"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MigrateDB(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(string(gooseDB.DialectSQLite3)); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func ResetDB(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(string(goose.DialectSQLite3)); err != nil {
		return err
	}

	if err := goose.Reset(db, "migrations"); err != nil {
		return err
	}
	return nil
}

type TxLike interface {
	DBTX
	Begin() (*sql.Tx, error)
}
