package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	JwtPropsKey string = "props"
)

func JwtValidation() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processJWT(w, r, next)
		})
	}
}

func processJWT(w http.ResponseWriter, r *http.Request, next http.Handler) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		errorValidation(w, http.StatusBadRequest, "", errors.New("malformed token"))
	} else {
		jwtToken := authHeader[1]
		claims, err := ValidateToken(jwtToken)
		if err != nil {
			errorValidation(w, http.StatusUnauthorized, "unauthorized", err)
		}
		ctx := context.WithValue(r.Context(), JwtPropsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func errorValidation(w http.ResponseWriter, status int, msg string, err error) {
	http.Error(w, errors.Wrap(err, msg).Error(), status)
}
