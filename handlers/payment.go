package handlers

import (
	"encoding/json"
	"merchant-bank-api/models"
	"merchant-bank-api/repository"
	"net/http"
	"time"
)

func Payment(w http.ResponseWriter, r *http.Request) {
	var req models.PaymentRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Validate amount
	if req.Amount <= 0 {
		http.Error(w, "Invalid amount. Must be greater than zero.", http.StatusBadRequest)
		return
	}

	customers, _ := repository.ReadCustomers()
	merchants, _ := repository.ReadMerchants()

	// Check if the customer is logged in
	for _, customer := range customers {
		if customer.ID == req.CustomerID && customer.LoggedIn {
			// Check if the merchant exists (assuming you provide merchant ID in the payment request)
			var merchantFound bool
			for _, merchant := range merchants {
				if merchant.ID == req.MerchantID {
					merchantFound = true
					break
				}
			}

			if !merchantFound {
				http.Error(w, "Invalid merchant ID", http.StatusBadRequest)
				return
			}

			// Create history log with the transaction amount and merchant ID
			history := models.History{
				CustomerID: req.CustomerID,
				MerchantID: req.MerchantID,
				Action:     "payment",
				Amount:     req.Amount,
				Timestamp:  time.Now().String(),
			}
			repository.AppendHistory(history)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Payment successful")
			return
		}
	}
	http.Error(w, "Payment failed. Customer not logged in or invalid ID", http.StatusUnauthorized)
}
