package db

import "time"

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   int64     `json:"balance" db:"balance"` // balance in cents
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
