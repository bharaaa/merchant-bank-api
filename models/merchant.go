package models

type Merchant struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"` // Amount the merchant can receive from transactions
}
