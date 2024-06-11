package middleware

import (
	"net/http"
	"strings"

	"ya-gophkeeper-server/pkg/jwt"

	"github.com/sirupsen/logrus"
)

func AuthMiddleware3(jwtKey []byte) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token, ok := strings.CutPrefix(authorization, "Bearer ")
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			userID, err := jwt.VerifyJWT(jwtKey, token)
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r.Header.Add("X-User-ID", userID)
			h.ServeHTTP(w, r)
		})
	}
}
