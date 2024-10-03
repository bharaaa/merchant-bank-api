package models

type History struct {
	CustomerID string  `json:"customer_id"`
	MerchantID string  `json:"merchant_id"`
	Action     string  `json:"action"`
	Amount     float64 `json:"amount,omitempty"` // Amount of the transactions
	Timestamp  string  `json:"timestamp"`
}
