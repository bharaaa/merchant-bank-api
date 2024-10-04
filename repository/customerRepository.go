package repository

import (
	"encoding/json"
	"merchant-bank-api/models"
	"os"
)

var customerFilePath = "data/customers.json"

func ReadCustomers() ([]models.Customer, error) {
	data, err := os.ReadFile(customerFilePath)
	if err != nil {
		return nil, err
	}
	var customers []models.Customer
	err = json.Unmarshal(data, &customers)
	return customers, err
}

func WriteCustomers(customers []models.Customer) error {
	data, err := json.MarshalIndent(customers, "", "  ")
	if err != nil {
		return err
	}

	// Write the updated data back to the file
	err = os.WriteFile(customerFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
