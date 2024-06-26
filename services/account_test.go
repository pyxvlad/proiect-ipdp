package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/rs/zerolog"
)

const email = "cat@meow.com"
const password = "catmeow"

func TestCreateAccount(t *testing.T) {
	t.Parallel()
	ctx := Context(t)

	FixtureAccount(ctx, t)
}

func FixtureAccount(ctx context.Context, t *testing.T) types.AccountID {
	t.Helper()
	return FixtureAccountWithSeed(ctx, t, "")
}

func FixtureAccountWithSeed(ctx context.Context, t *testing.T, seed string) types.AccountID {
	t.Helper()

	log := zerolog.Ctx(ctx).With().Str("test_name", t.Name()).Logger()

	fixtureAccounts := ctx.Value(ContextKeyFixtureAccountsCache).(FixtureAccountsCache)

	accountID, found := fixtureAccounts[seed]
	if found {
		log.Debug().
			Int64("account_id", int64(accountID)).
			Msg("found account in fixture account cache")
		return accountID
	}

	as := services.NewAccountService()
	accountID, err := as.CreateAccountWithEmail(ctx, services.AccountData{
		Email:    seed + email,
		Password: password,
	})

	if err != nil {
		t.Fatal(err)
	}

	fixtureAccounts[seed] = accountID;

	log.Debug().
		Int64("account_id", int64(accountID)).
		Msg("freshly made")

	return accountID
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
	t.Parallel()
	ctx := Context(t)

	FixtureAccount(ctx, t)

	as := services.NewAccountService()
	_, err := as.CreateAccountWithEmail(ctx, services.AccountData{
		Email:    email,
		Password: password,
	})

	t.Logf("%#v", err)
	if !errors.Is(err, nil) {
		t.Fatal("Should have gotten a duplicate key for the accounts.email")
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	ctx := Context(t)

	FixtureAccount(ctx, t)
	FixtureSession(ctx, t)
}

func TestGetAccountFromSession(t *testing.T) {
	t.Parallel()
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
