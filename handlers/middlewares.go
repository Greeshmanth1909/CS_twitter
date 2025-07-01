package handlers

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

var jwtClaims ctxKey = "claims"

func AuthMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Header.Get("Authorization")
		jwt = strings.TrimPrefix(jwt, "Bearer ")

		jwtclaims, err := verifyToken(jwt)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Please Login to continue"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, jwtClaims, jwtclaims)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
