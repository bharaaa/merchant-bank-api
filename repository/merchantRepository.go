package repository

import (
	"encoding/json"
	"merchant-bank-api/models"
	"os"
)

var merchantFilePath = "data/merchants.json"

func ReadMerchants() ([]models.Merchant, error) {
	data, err := os.ReadFile(merchantFilePath)
	if err != nil {
		return nil, err
	}

	var merchants []models.Merchant
	if err := json.Unmarshal(data, &merchants); err != nil {
		return nil, err
	}
	return merchants, nil
}

func WriteMerchants(merchants []models.Merchant) error {
	// data, err := json.Marshal(merchants)
	// if err != nil {
	// 	return err
	// }
	// return os.WriteFile("data/merchants.json", data, 0644)

	data, err := json.MarshalIndent(merchants, "", "  ")
	if err != nil {
		return err
	}

	// Write the updated data back to the file
	err = os.WriteFile(merchantFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AppendMerchant(merchant models.Merchant) error {
	merchants, err := ReadMerchants()
	if err != nil {
		return err
	}

	merchants = append(merchants, merchant)
	return WriteMerchants(merchants)
}
