package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const CtxUserID ctxKey = "user_id"

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				unauth(w)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			tkn, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err != nil || !tkn.Valid {
				unauth(w)
				return
			}

			claims, ok := tkn.Claims.(jwt.MapClaims)
			if !ok {
				unauth(w)
				return
			}
			uidFloat, ok := claims["user_id"].(float64)
			if !ok {
				unauth(w)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, int(uidFloat))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauth(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"unauthorized"}`))
}
