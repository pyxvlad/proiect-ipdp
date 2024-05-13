package routes

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pyxvlad/proiect-ipdp/handlers"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/rs/zerolog"
)

func NewAppRouter(log *zerolog.Logger, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.WithContext(r.Context())
			ctx = context.WithValue(ctx, services.ContextKeyDB, db)
			requestWithLogger := r.WithContext(ctx)
			next.ServeHTTP(w, requestWithLogger)
		})
	})

	hello := handlers.NewHelloHandler(
		log.With().Str("handler", "hello-handler").Logger(),
	)
	router.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	router.Method(http.MethodGet, "/hello", hello)
	router.Route("/signup", func(r chi.Router) {
		r.Get("/", handlers.SignUpPage)
		r.Post("/attempt", handlers.SignUpAttempt)
	})

	router.Route("/login", func(r chi.Router) {
		r.Get("/", handlers.LogInPage)
		r.Post("/attempt", handlers.LogInAttempt)
	})

	router.HandleFunc("/samples", handlers.SampleBookCards)
	router.HandleFunc("/addbook", handlers.AddBookPage)
	router.Post("/books/cards/preview", handlers.PreviewCard)

	fs := http.FileServer(http.Dir("./assets/"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	return router
}

func shutdownHandler(server *http.Server, log *zerolog.Logger) {
	const sigintChannelSize = 1
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

func ListenAndServe(log *zerolog.Logger, db *sql.DB) {
	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = NewAppRouter(log, db)
	go shutdownHandler(server, log)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("error while shutting down server")
	}
}
