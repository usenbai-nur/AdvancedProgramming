package auth

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const usernameKey ctxKey = "username"

func UsernameFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(usernameKey)
	s, ok := v.(string)
	return s, ok
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}

		// Expect: "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		username, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), usernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
