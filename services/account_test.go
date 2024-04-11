package services_test

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/pyxvlad/proiect-ipdp/models"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func FreshDB(t *testing.T) *gorm.DB {
	t.Helper()

	dbPath := path.Join(t.TempDir(), "tmp.db")
	sqliteDB := sqlite.Open(dbPath)

	db, err := gorm.Open(sqliteDB, &gorm.Config{TranslateError: true})

	if err != nil {
		t.Fatal(err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		t.Fatal(err)
	}

	return db
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

	ctx := context.Background()
	ctx = log.WithContext(ctx)
	ctx = context.WithValue(ctx, services.ContextKeyDB, db)

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

func FixtureSession(ctx context.Context, t *testing.T) models.Session {
	t.Helper()
	as := services.NewAccountService()
	account, err := as.Login(ctx, services.AccountData{
		Email:    email,
		Password: password,
	})

	if err != nil {
		t.Fatal(err)
	}

	session, err := as.CreateSession(ctx, account.ID)
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
	if !errors.Is(err, gorm.ErrDuplicatedKey) {
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
	session := FixtureSession(ctx, t)
	t.Logf("session: %#v\n", session)

	as := services.NewAccountService()
	account, err := as.GetAccountForSession(ctx, session.Token)
	if err != nil {
		t.Fatal(err)
	}

	if session.AccountID != account.ID {
		t.Fatal("I got the wrong account from the token")
	}
}
