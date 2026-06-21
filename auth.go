package main

import (
	"context"
	"net/http"
)

type ctxUserKey struct{}

func GetUserFromContext(r *http.Request) *User {
	v := r.Context().Value(ctxUserKey{})
	if v == nil {
		return nil
	}
	u, _ := v.(*User)
	return u
}

// AuthMiddleware validates the token provided in `X-User-Token` header and
// injects the user into the request context.
func (a *App) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-User-Token")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		user, err := a.Store.GetUserByToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserKey{}, user)
		next(w, r.WithContext(ctx))
	}
}

// RequireRole wraps a handler and ensures the authenticated user has the given role.
func (a *App) RequireRole(role Role, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if user.Role != role {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
