package services

import (
	"context"

	"github.com/pyxvlad/proiect-ipdp/models"
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
) error {

	log := zerolog.Ctx(ctx).With().Caller().Str("email", request.Email).Logger()
	log.Info().Msg("creating account")

	hashed, err := service.hasher(request.Password)
	if err != nil {
		return err
	}

	result := DB(ctx).Create(&models.Account{
		Email:    request.Email,
		Password: hashed,
	})

	if result.Error != nil {
		log.Err(result.Error).Send()
		return result.Error
	}

	log.Info().Msg("created account")
	return nil
}

func (service *AccountService) Login(ctx context.Context, loginData AccountData) (models.Account, error) {
	log := zerolog.Ctx(ctx).With().Caller().Str("email", loginData.Email).Logger()
	log.Info().Msg("logging in account")

	var account models.Account

	result := DB(ctx).
		Where("email = ?", loginData.Email).
		Take(&account)

	if result.Error != nil {
		log.Err(result.Error).Msg("failed to login")
		return models.Account{}, result.Error
	}

	err := service.comparator(account.Password, loginData.Password)
	if err != nil {
		log.Info().Msg("account had wrong password")
		return models.Account{}, err
	}

	log.Info().Msg("account logged in")
	return account, nil
}
