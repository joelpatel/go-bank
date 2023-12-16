package db

import "time"

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   int64     `json:"balance" db:"balance"` // balance in cents
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Entry struct {
	ID        int64     `json:"id" db:"id"`
	AccountID int64     `json:"account_id" db:"account_id"`
	Amount    int64     `json:"amount" db:"amount"` // amount in cents
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Transfer struct {
	ID            int64     `json:"id" db:"id"`
	FromAccountID int64     `json:"from_account_id" db:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id" db:"to_account_id"`
	Amount        int64     `json:"amount" db:"amount"` // amount in cents
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type TransferTxResult struct {
	TransferRecord  Transfer `json:"transfer"`
	FromAccount     Account  `json:"from_account"`
	ToAccount       Account  `json:"to_account"`
	FromEntryRecord Entry    `json:"from_entry"`
	ToEntryRecord   Entry    `json:"to_entry"`
}
