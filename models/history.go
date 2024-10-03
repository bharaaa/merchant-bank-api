package models

type History struct {
	CustomerID string  `json:"customer_id"`
	Action     string  `json:"action"`
	Amount     float64 `json:"amount,omitempty"`
	Timestamp  string  `json:"timestamp"`
}
