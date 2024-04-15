package services_test

import (
	"context"
	"database/sql"
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/rs/zerolog"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func FreshDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := path.Join(t.TempDir(), "tmp.db")
	sqliteDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}

	err = database.MigrateDB(sqliteDB)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sqliteDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatal("Failed to enable foreign keys:", err)
	}

	return sqliteDB
}

func FreshLog(t *testing.T) *zerolog.Logger {
	t.Helper()

	logPath := path.Join(t.TempDir(), "test.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		err := logFile.Close()
		if err != nil {
			panic(err)
		}
	})
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	writer := zerolog.MultiLevelWriter(logFile, consoleWriter)
	log := zerolog.New(writer).With().Timestamp().Logger()
	log.Debug().Msgf("log path: %s", logPath)
	return &log
}

func Context(t *testing.T) context.Context {
	t.Helper()
	log := FreshLog(t)
	db := FreshDB(t)

	var dbtx database.TxLike = db
	ctx := context.Background()
	ctx = log.WithContext(ctx)
	ctx = context.WithValue(ctx, services.ContextKeyDB, dbtx)

	return ctx
}
