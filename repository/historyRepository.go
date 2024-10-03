package repository

import (
	"encoding/json"
	"merchant-bank-api/models"
	"os"
)

func AppendHistory(history models.History) error {
	data, err := os.ReadFile("data/history.json")
	if err != nil {
		return err
	}
	var histories []models.History
	json.Unmarshal(data, &histories)
	histories = append(histories, history)
	newData, _ := json.Marshal(histories)
	return os.WriteFile("data/history.json", newData, 0644)
}
