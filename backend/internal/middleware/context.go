package middleware

import (
	"context"
	"net/http"

	"github.com/fajrinajiseno/mygolangapp/internal/config"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Context().Value(config.ContextUserID); v != nil {
			next.ServeHTTP(w, r)
			return
		}

		sub, err := GetTokenSub(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		newReq := r.WithContext(context.WithValue(r.Context(), config.ContextUserID, sub))
		next.ServeHTTP(w, newReq)
	})
}
