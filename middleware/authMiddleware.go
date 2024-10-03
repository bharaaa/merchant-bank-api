package middleware

import (
	"encoding/json"
	"merchant-bank-api/config"
	"merchant-bank-api/handlers"
	"merchant-bank-api/tokenblacklist"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Parse the token from the Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if the token is blacklisted
		if tokenblacklist.IsTokenBlacklisted(tokenString) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Token is invalid"})
			return
		}

		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Signature"})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token"})
			return
		}

		if !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token"})
			return
		}

		// Check token expiration
		if claims.ExpiresAt.Before(time.Now()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Token expired"})
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
