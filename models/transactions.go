package models

import "time"

type Transaction struct {
	ID        int64     `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // "credit" or "debit"
	CreatedAt time.Time `json:"created_at"`
}
