package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func jwtAuthHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			return []byte("123456789"), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}
