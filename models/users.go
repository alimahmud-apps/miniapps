package models

import "time"

// User represents a user in the database.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
