package handlers

import (
	"encoding/json"
	"merchant-bank-api/config"
	"merchant-bank-api/models"
	"merchant-bank-api/repository"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	CustomerID string `json:"customer_id"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var reqCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&reqCustomer); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	customers, err := repository.ReadCustomers()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Error reading customers data"}`, http.StatusInternalServerError)
		return
	}

	for _, customer := range customers {
		if customer.Username == reqCustomer.Username && customer.Password == reqCustomer.Password {
			// Create JWT claims (payload)
			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &Claims{
				CustomerID: customer.ID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
				},
			}

			// Create the token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(config.JwtKey)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error": "Could not create token"}`, http.StatusInternalServerError)
				return
			}

			// Respond with the JWT
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "Customer not found or invalid credentials"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Set token cookie to expire
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
