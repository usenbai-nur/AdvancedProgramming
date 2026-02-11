package auth

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const (
	usernameKey ctxKey = "username"
	roleKey     ctxKey = "role"
)

func UsernameFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(usernameKey)
	s, ok := v.(string)
	return s, ok
}

func RoleFromContext(ctx context.Context) (Role, bool) {
	v := ctx.Value(roleKey)
	r, ok := v.(Role)
	return r, ok
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), usernameKey, claims.Username)
		ctx = context.WithValue(ctx, roleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRoles(next http.Handler, roles ...Role) http.Handler {
	allowed := make(map[Role]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := RoleFromContext(r.Context())
		if !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if _, ok := allowed[role]; !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}))
}
