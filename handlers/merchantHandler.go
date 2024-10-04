package handlers

import (
	"encoding/json"
	"merchant-bank-api/models"
	"merchant-bank-api/repository"
	"net/http"

	"github.com/google/uuid"
)

func CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var newMerchant models.Merchant

	// Decode the request body into the newMerchant struct
	err := json.NewDecoder(r.Body).Decode(&newMerchant)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}

	// Generate new UUID for the merchant
	newMerchant.ID = uuid.New().String()
	newMerchant.Balance = 0.0

	// Read existing merchants from the JSON file
	merchants, err := repository.ReadMerchants()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error reading merchant data"})
		return
	}

	// Append the new merchant to the merchants slice
	merchants = append(merchants, newMerchant)

	// Write the updated merchants list to the JSON file
	err = repository.WriteMerchants(merchants)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error writing merchant data"})
		return
	}

	// Respond with the newly created merchant (with UUID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": newMerchant.ID, "name": newMerchant.Name})
}
