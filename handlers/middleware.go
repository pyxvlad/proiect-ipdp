package handlers

import (
	"context"
	"net/http"

	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/rs/zerolog"
)

type ContextKey string

const (
	ContextKeyAccountID ContextKey = "ck_account_id"
)

// The ContextKeyAccountID key must have a value associated with it inside the ctx.
func AccountID(ctx context.Context) types.AccountID {
	//nolint:revive // See above.
	return ctx.Value(ContextKeyAccountID).(types.AccountID)
}

func LoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)

		cookie, err := r.Cookie("token")
		if err != nil {
			w.Header().Set("Location", "/login")
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(http.StatusFound)
			return
		}

		accountService := services.NewAccountService()
		accountID, err := accountService.GetAccountForSession(ctx, cookie.Value)

		if err != nil {
			w.Header().Set("Location", "/login")
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(http.StatusUnauthorized)
			log.Info().Err(err).Str("token", cookie.Value).Msg("while trying to validate token")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, ContextKeyAccountID, accountID)))
	})
}
