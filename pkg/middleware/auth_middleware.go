package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		autheHandler := r.Header.Get("Authorization")
		token := strings.TrimPrefix(autheHandler, "Bearer ")
		fmt.Println(token)
		next.ServeHTTP(w, r)
	})
}
