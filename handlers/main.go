package handlers

import (
	"net/http"

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
