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
	json.NewDecoder(r.Body).Decode(&reqCustomer)

	customers, _ := repository.ReadCustomers()

	for _, customer := range customers {
		if customer.Username == reqCustomer.Username && customer.Password == reqCustomer.Password {
			// Set customer to logged in
			customer.LoggedIn = true
			repository.WriteCustomers(customers)

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
				http.Error(w, "Could not create token", http.StatusInternalServerError)
				return
			}

			// Respond with the JWT
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
			return
		}
	}
	http.Error(w, "Customer not found or invalid credentials", http.StatusUnauthorized)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Set token cookie to expire
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}
