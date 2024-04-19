/*
 * This file was last modified at 2024-04-20 00:25 by Victor N. Skurikhin.
 * jwt.go
 * $Id$
 */

package utils

import (
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/model"
	"net/http"
)

const Secret = "<jwt-secret>" // Replace <jwt-secret> with your secret key that is private to you.

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
}

func MakeToken(name string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"username": name})
	return tokenString
}

func Verifier() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verifier(tokenAuth)(next)
	}
}

func UnauthenticatedError(next http.Handler) http.Handler {
	return unAuthError(next, handlers.ErrUserUnauthenticated)
}

func UnauthorizedError(next http.Handler) http.Handler {
	return unAuthError(next, handlers.ErrUserUnauthorized)
}

func unAuthError(next http.Handler, err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := jwtauth.FromContext(r.Context())

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, "", http.StatusUnauthorized)
			//goland:noinspection GoUnhandledErrorResult
			render.Render(w, r, model.Error(err))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
