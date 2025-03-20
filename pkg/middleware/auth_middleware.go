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
	ContextUserKey key = "ContextUserKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuth(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearerToken, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		token := strings.TrimPrefix(bearerToken, "Bearer ")

		data, err := jwt.NewJWT(cfg.Auth.Secret).ParseToken(token)
		if err != nil {
			writeUnauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, data)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
