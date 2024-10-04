package handlers

import (
	"encoding/json"
	"merchant-bank-api/models"
	"merchant-bank-api/repository"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer models.Customer

	// Decode the request body into the newCustomer struct
	err := json.NewDecoder(r.Body).Decode(&newCustomer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}

	// Ensure both username and password are provided
	if newCustomer.Username == "" || newCustomer.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Username and password are required"})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newCustomer.Password), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error hashing password"})
		return
	}

	// Set the hashed password
	newCustomer.Password = string(hashedPassword)

	// Generate a new UUID for the customer
	newCustomer.ID = uuid.New().String()

	// Read existing customers from the JSON file
	customers, err := repository.ReadCustomers()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error reading customer data"})
		return
	}

	// Append the new customer to the customers slice
	customers = append(customers, newCustomer)

	// Write the updated customers list to the JSON file
	err = repository.WriteCustomers(customers)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error writing customer data"})
		return
	}

	// Respond with the newly created customer (with UUID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": newCustomer.ID, "username": newCustomer.Username})
}
