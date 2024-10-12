package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/esafronov/yp-sysloyalty/internal/helpers/token"
)

func JwtAuthMiddleware(secret string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := w
			authHeader := r.Header.Get("Authorization")
			t := strings.Split(authHeader, " ")
			if len(t) == 2 {
				authToken := t[1]
				authorized, _ := token.IsAuthorized(authToken, secret)
				if authorized {
					userID, err := token.ExtractIDFromToken(authToken, secret)
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					nr := r.WithContext(context.WithValue(r.Context(), "x-user-id", userID))
					h.ServeHTTP(sw, nr)
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
