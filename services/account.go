package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHashGenerator func(password string) (string, error)
type PasswordHashComparator func(hashed string, password string) error

const BcryptWorkFactor = 12

func generateBcryptPasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		BcryptWorkFactor,
	)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func compareBcryptPasswordHash(hashed string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

type AccountService struct {
	hasher     PasswordHashGenerator
	comparator PasswordHashComparator
}

func NewAccountService() AccountService {
	return AccountService{
		hasher:     generateBcryptPasswordHash,
		comparator: compareBcryptPasswordHash,
	}
}

type AccountData struct {
	// User's email
	Email string `json:"email"`
	// User's password
	Password string `json:"password"`
}

func (service *AccountService) CreateAccountWithEmail(
	ctx context.Context, request AccountData,
) (types.AccountID, error) {
	log := zerolog.Ctx(ctx).With().Caller().Str("email", request.Email).Logger()
	log.Info().Msg("creating account")

	hashed, err := service.hasher(request.Password)
	if err != nil {
		return types.InvalidAccountID, err
	}

	accountID, err := DB(ctx).CreateAccountWithEmail(ctx, database.CreateAccountWithEmailParams{
		Email:    request.Email,
		Password: hashed,
	})

	if err != nil {
		log.Err(err).Send()
		return types.InvalidAccountID, err
	}

	log.Info().Msg("created account")
	return accountID, nil
}

func (service *AccountService) Login(
	ctx context.Context, loginData AccountData,
) (types.AccountID, error) {
	log := zerolog.Ctx(ctx).
		With().
		Caller().
		Str("email", loginData.Email).
		Logger()

	log.Info().Msg("logging in account")

	row, err := DB(ctx).GetPasswordByEmail(ctx, loginData.Email)

	if err != nil {
		log.Err(err).Msg("failed to login")
		return 0, err
	}

	err = service.comparator(row.Password, loginData.Password)
	if err != nil {
		log.Info().Msg("account had wrong password")
		return 0, err
	}

	log.Info().Msg("account logged in")
	return row.AccountID, nil
}

func (service *AccountService) CreateSession(
	ctx context.Context, accountID types.AccountID,
) (string, error) {
	log := zerolog.Ctx(ctx).
		With().
		Caller().
		Logger()

	log.Info().Msg("creating session")

	const TOKEN_BYTES = 32
	var token [TOKEN_BYTES]byte
	count, err := rand.Read(token[:])
	if err != nil {
		log.Err(err).Send()
		return "", err
	}

	if count < TOKEN_BYTES {
		log.Error().Msgf("read %d/%d out of required bytes for a token from crypto/rand", count, TOKEN_BYTES)
	}

	var sessionData database.CreateSessionTokenParams
	sessionData.Token = base64.RawStdEncoding.EncodeToString(token[:])
	sessionData.AccountID = accountID

	err = DB(ctx).CreateSessionToken(ctx, sessionData)

	if err != nil {
		return "", err
	}

	return sessionData.Token, nil
}

func (as *AccountService) GetAccountForSession(
	ctx context.Context, token string,
) (types.AccountID, error) {
	accountID, err := DB(ctx).GetSessionAccount(ctx, token)

	return accountID, err
}
