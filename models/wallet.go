package models

type Wallet struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Balance      float64  `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID     string  `json:"id"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}