package entity

import "time"

type Payment struct {
	ID        string    `json:"id"`
	Merchant  string    `json:"merchant"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentSummary struct {
	TotalByFiler   int
	Total          int
	TotalCompleted int
	TotalFailed    int
	TotalPending   int
}
