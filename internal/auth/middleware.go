package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/soniceroy/hackernewsGQL/internal/users"
	"github.com/soniceroy/hackernewsGQL/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Middleware jwt auth middleware
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate jwt token
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// check if user exists in db
			user := users.User{Username: username}
			id, err := users.GetUserIDByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.ID = strconv.Itoa(id)
			// put user in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user in the context. REQUIRES Middleware to run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
