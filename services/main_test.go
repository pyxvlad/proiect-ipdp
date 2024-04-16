package services_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/rs/zerolog"
)

var databasePath string

func TestMain(m *testing.M) {

	tmpdir, err := os.MkdirTemp("", "ipdp-test-db-*")
	if err != nil {
		panic(err)
	}

	dbPath := path.Join(tmpdir, "template.db")
	sqliteDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	err = database.MigrateDB(sqliteDB)
	if err != nil {
		panic(err)
	}

	_, err = sqliteDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		panic(err)
	}

	err = sqliteDB.Close()
	if err != nil {
		panic(err)
	}

	databasePath = tmpdir


	exitCode := m.Run()

	if exitCode == 0 {
		os.RemoveAll(tmpdir)
	} else {
		fmt.Println("TEST: the path with the databases is:", databasePath)
	}

	os.Exit(exitCode)
}

func FreshDB(t *testing.T) *sql.DB {
	t.Helper()

	mapper := func(r rune) rune {
		if r < utf8.RuneSelf {
			const allowed = "!#$%&()+,-.=@^_{}~ "
			if '0' <= r && r <= '9' ||
				'a' <= r && r <= 'z' ||
				'A' <= r && r <= 'Z' {
				return r
			}
			if strings.ContainsRune(allowed, r) {
				return r
			}
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}
	pattern := strings.Map(mapper, t.Name())

	templateDB, err := os.ReadFile(path.Join(databasePath, "template.db"))
	if err != nil {
		t.Fatal(err)
	}
	dbFile := path.Join(databasePath, pattern)
	err = os.WriteFile(dbFile, templateDB, 0666)
	if err != nil {
		t.Fatal(err)
	}

	sqliteDB, err := sql.Open("sqlite3", dbFile)
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
