package routes

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pyxvlad/proiect-ipdp/handlers"
	"github.com/rs/zerolog"
)

func NewAppRouter(log *zerolog.Logger) *chi.Mux {
	router := chi.NewRouter()
	hello := handlers.NewHelloHandler(
		log.With().Str("handler", "hello-handler").Logger(),
	)
	router.Method(http.MethodGet, "/hello", hello)

	return router
}

func shutdownHandler(server *http.Server, log *zerolog.Logger) {
	const sigintChannelSize = 1;
	sigint := make(chan os.Signal, sigintChannelSize)
	defer close(sigint)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	const shutdownTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Err(err).Msg("shutdown with error")
	}
}

func ListenAndServe(log *zerolog.Logger) {
	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = NewAppRouter(log)
	go shutdownHandler(server, log)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("error while shutting down server")
	}
}
