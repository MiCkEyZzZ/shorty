package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(mid ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mid) - 1; i >= 0; i-- {
			next = mid[i](next)
		}
		return next
	}
}
