package handlers

import (
	"net/http"
	"net/mail"

	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/pyxvlad/proiect-ipdp/templates"
	"github.com/rs/zerolog"
)

type HelloHandler struct {
	log zerolog.Logger
}

func NewHelloHandler(log zerolog.Logger) http.Handler {
	handler := new(HelloHandler)

	handler.log = log

	return handler
}

// ServeHTTP implements http.Handler.
func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := templates.HelloPage().Render(r.Context(), w)
	h.log.Err(err).Send()
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := templates.SignUpPage().Render(r.Context(), w)
	log.Err(err).Send()
}

func SignUpAttempt(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())

	err := r.ParseForm()
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	_, err = mail.ParseAddress(email)
	if err != nil {
		err = templates.SignUpForm("Invalid email address, retry.").Render(r.Context(), w)
		log.Err(err).Send()
		return
	}

	if password != confirmPassword {
		err = templates.SignUpForm("Passwords don't match, retry.").Render(r.Context(), w)
		log.Err(err).Send()
		return
	}

	as := services.NewAccountService()
	as.CreateAccountWithEmail(r.Context(), services.AccountData{
		Email:    email,
		Password: password,
	})

	w.Header().Add("HX-Location", "/hello")

	w.WriteHeader(http.StatusCreated)
}

func LogInPage(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := templates.LogInPage().Render(r.Context(), w)
	log.Err(err).Send()
}

func LogInAttempt(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	accountService := services.NewAccountService()
	account, err := accountService.Login(r.Context(), services.AccountData{
		Email:    email,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	session, err := accountService.CreateSession(r.Context(), account.ID)
	if err != nil {
		panic(err)
	}
	cookie := http.Cookie{
		Name:   "token",
		Value:  session.Token,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/hello")
	_ = log
}

func AddBookPage(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := templates.AddBookPage().Render(r.Context(), w)
	log.Err(err).Send()
}
