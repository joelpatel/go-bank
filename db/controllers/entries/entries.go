package entries

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/joelpatel/go-bank/db"
)

const (
	INSERT_ENTRY_QUERY = "INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING id, account_id, amount, created_at;"
)

// scan entry row
func ScanRow(row *sql.Row) (*Entry, error) {
	var entry Entry

	err := row.Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// create (self contained transaction)
func CreateEntryTx(accountID, amount int64) (*Entry, error) {
	tx := db.Conn.MustBegin()

	row := tx.QueryRow(INSERT_ENTRY_QUERY, accountID, amount)

	entry, err := ScanRow(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return entry, nil
}

// create
func CreateEntry(tx *sqlx.Tx, accountID, amount int64) (*Entry, error) {
	return ScanRow(tx.QueryRow(INSERT_ENTRY_QUERY, accountID, amount))
}

// read
func GetEntryByID(id string) (*Entry, error) {
	var entry Entry

	err := db.Conn.Get(&entry, "SELECT id, account_id, amount, created_at FROM entries WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// read all for account_id (pagination)
func GetEntriesByAccountID(account_id, limit, offset int64) (*[]Entry, error) {
	var entries []Entry

	err := db.Conn.Select(&entries, "SELECT id, account_id, amount, created_at FROM entries WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3;", account_id, limit, offset)
	if err != nil {
		return nil, err
	}

	return &entries, nil
}
