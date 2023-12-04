package entries

import "time"

type Entry struct {
	ID        int64     `json:"id" db:"id"`
	AccountID int64     `json:"account_id" db:"account_id"`
	Amount    int64     `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
