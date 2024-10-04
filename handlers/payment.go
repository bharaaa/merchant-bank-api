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
)

func Payment(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Parse token from the Authorization header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Check if token is blacklisted
	if tokenblacklist.IsTokenBlacklisted(tokenString) {
		http.Error(w, "Token is invalid", http.StatusUnauthorized)
		return
	}

	// Parse the token and get the claims
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Token is invalid", http.StatusUnauthorized)
		return
	}

	// Extract the customer ID from JWT claims
	customerIDFromToken := claims.CustomerID

	var req models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid amount. Must be greater than zero."})
		return
	}

	customers, err := repository.ReadCustomers()
	if err != nil {
		http.Error(w, "Error reading customer data", http.StatusInternalServerError)
		return
	}

	merchants, err := repository.ReadMerchants()
	if err != nil {
		http.Error(w, "Error reading merchant data", http.StatusInternalServerError)
		return
	}

	// Verify the customer exists and matches the customer ID from the token
	var customerFound bool
	for _, customer := range customers {
		if customer.ID == customerIDFromToken {
			customerFound = true
			break
		}
	}

	if !customerFound {
		http.Error(w, "Payment failed. Customer not found or not logged in", http.StatusUnauthorized)
		return
	}

	// Verify if the merchant exists
	var merchantFound bool
	for i, merchant := range merchants {
		if merchant.ID == req.MerchantID {
			merchantFound = true
			merchants[i].Balance += req.Amount
			break
		}
	}

	if !merchantFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid merchant ID"})
		return
	}

	// Write the updated merchant data back to the merchants.json file
	err = repository.WriteMerchants(merchants)
	if err != nil {
		http.Error(w, "Error updating merchant data", http.StatusInternalServerError)
		return
	}

	// Log the payment history
	history := models.History{
		CustomerID: customerIDFromToken,
		MerchantID: req.MerchantID,
		Action:     "payment",
		Amount:     req.Amount,
		Timestamp:  time.Now().String(),
	}

	err = repository.AppendHistory(history)
	if err != nil {
		http.Error(w, "Failed to log payment history", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Payment successful"})
}
