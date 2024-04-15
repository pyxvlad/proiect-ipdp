package services_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/rs/zerolog"
)

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

const email = "cat@meow.com"
const password = "catmeow"

func TestCreateAccount(t *testing.T) {
	ctx := Context(t)

	FixtureAccount(ctx, t)
}

func FixtureAccount(ctx context.Context, t *testing.T) {
	as := services.NewAccountService()
	err := as.CreateAccountWithEmail(ctx, services.AccountData{
		Email:    email,
		Password: password,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func FixtureSession(ctx context.Context, t *testing.T) string {
	t.Helper()
	as := services.NewAccountService()
	accountID, err := as.Login(ctx, services.AccountData{
		Email:    email,
		Password: password,
	})

	if err != nil {
		t.Fatal(err)
	}

	session, err := as.CreateSession(ctx, accountID)
	if err != nil {
		t.Fatal(err)
	}

	return session
}

func TestDuplicateAccount(t *testing.T) {
	ctx := Context(t)

	FixtureAccount(ctx, t)

	as := services.NewAccountService()
	err := as.CreateAccountWithEmail(ctx, services.AccountData{
		Email:    email,
		Password: password,
	})

	t.Logf("%#v", err)
	if !errors.Is(err, nil) {
		t.Fatal("Should have gotten a duplicate key for the accounts.email")
	}
}

func TestLogin(t *testing.T) {
	ctx := Context(t)

	FixtureAccount(ctx, t)

	as := services.NewAccountService()
	_, err := as.Login(ctx, services.AccountData{
		Email:    email,
		Password: password,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateSession(t *testing.T) {
	ctx := Context(t)

	FixtureAccount(ctx, t)
	FixtureSession(ctx, t)
}

func TestGetAccountFromSession(t *testing.T) {
	ctx := Context(t)

	FixtureAccount(ctx, t)

	accountFromEmail, err := services.DB(ctx).GetAccountByEmail(ctx, email)
	if err != nil {
		t.Fatal(err)
	}

	token := FixtureSession(ctx, t)
	t.Logf("session: %#v\n", token)

	as := services.NewAccountService()
	accountFromToken, err := as.GetAccountForSession(ctx, token)
	if err != nil {
		t.Fatal(err)
	}

	if accountFromEmail != accountFromToken {
		t.Fatal("I got the wrong account from the token")
	}
}
