package middleware

import (
	"context"
	"net/http"
	"strings"

	"shorty/internal/config"
	"shorty/pkg/jwt"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuth(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHandler := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHandler, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		token := strings.TrimPrefix(authHandler, "Bearer ")
		isValid, data := jwt.NewJWT(cfg.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
