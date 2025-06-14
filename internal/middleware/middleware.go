package middleware

import (
	"context"
	"fmt"
	"goBackendServer/internal/store"
	"goBackendServer/internal/tokens"
	"goBackendServer/internal/utils"
	"net/http"
	"strings"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user in request") // bad actor call
	}
	return user
}

func (middleware *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Anon function to interject any incoming requests to our server
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := middleware.UserStore.GetUserToken(tokens.ScopeAuth, token)

		if err != nil {
			fmt.Println("%w", err)
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this feature"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
