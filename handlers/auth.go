package handlers

import (
	"encoding/json"
	"merchant-bank-api/config"
	"merchant-bank-api/models"
	"merchant-bank-api/repository"
	"merchant-bank-api/tokenblacklist"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	CustomerID string `json:"customer_id"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var reqCustomer models.Customer

	// Decode the request body into reqCustomer struct
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

	// Iterate through the customers to find matching username
	for _, customer := range customers {
		if customer.Username == reqCustomer.Username {
			// Compare the stored hashed password with the provided password
			err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(reqCustomer.Password))
			if err != nil {
				// Password doesn't match
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
				return
			}

			// Password matches, proceed with generating JWT
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

	// If no matching customer found
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
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   "",
	// 	Expires: time.Now().Add(-time.Hour),
	// })

	// Parse the token from the Authorization header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Add the token to the blacklist
	tokenblacklist.AddToken(tokenString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
