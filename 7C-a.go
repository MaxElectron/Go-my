//go:build !solution

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type User struct {
	Name  string
	Email string
}

func ContextUser(ctx context.Context) (*User, bool) {
	context_value := ctx.Value(0)
	return context_value.(*User), context_value != nil
}

var ErrInvalidToken = errors.New("invalid token")

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
				head := r.Header.Get("Authorization")[7:]
				if u, err := checker.CheckToken(r.Context(), head); err != nil {
					wr.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					context := context.WithValue(r.Context(), 0, u)
					next.ServeHTTP(wr, r.WithContext(context))
				}
			} else {
				wr.WriteHeader(http.StatusUnauthorized)
				return
			}

		})
	}
}
