package repository

import (
	"encoding/json"
	"merchant-bank-api/models"
	"os"
)

func ReadCustomers() ([]models.Customer, error) {
	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		return nil, err
	}
	var customers []models.Customer
	err = json.Unmarshal(data, &customers)
	return customers, err
}

func WriteCustomers(customers []models.Customer) error {
	data, err := json.Marshal(customers)
	if err != nil {
		return err
	}
	return os.WriteFile("data/customers.json", data, 0644)
}
