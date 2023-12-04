package transfers

import "time"

type Transfer struct {
	ID            int64     `json:"id" db:"id"`
	FromAccountID int64     `json:"from_account_id" db:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id" db:"to_account_id"`
	Amount        int64     `json:"amount" db:"amount"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
