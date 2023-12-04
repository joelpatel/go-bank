package accounts

import "time"

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   int64     `json:"balance" db:"balance"` // int64 scaled representing last 3 digits as decimal points
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
