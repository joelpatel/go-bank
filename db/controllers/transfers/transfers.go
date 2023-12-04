package transfers

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/joelpatel/go-bank/db"
)

const (
	INSERT_TRANSFER_QUERY = "INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3) RETURNING id, from_account_id, to_account_id, amount, created_at;"
)

func ScanRow(row *sql.Row) (*Transfer, error) {
	var transfer Transfer

	err := row.Scan(&transfer.ID, &transfer.FromAccountID, &transfer.ToAccountID, &transfer.Amount, &transfer.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// create (self contained transaction)
func CreateTransferTx(from_account_id, to_account_id, amount int64) (*Transfer, error) {
	tx := db.Conn.MustBegin()

	row := tx.QueryRow(INSERT_TRANSFER_QUERY, from_account_id, to_account_id, amount)

	transfer, err := ScanRow(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return transfer, nil
}

// create
func CreateTransfer(tx *sqlx.Tx, from_account_id, to_account_id, amount int64) (*Transfer, error) {
	return ScanRow(tx.QueryRow(INSERT_TRANSFER_QUERY, from_account_id, to_account_id, amount))
}

// read (id)
func GetTransferByID(id string) (*Transfer, error) {
	var transfer Transfer

	err := db.Conn.Get(&transfer, "SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// read (from_account_id OR to_account_id)
// (-1 if don't want to search for from exor to)
// paginated
func GetTransfersFromTo(from_account_id, to_account_id, limit, offset int64) (*[]Transfer, error) {
	var transfers []Transfer

	err := db.Conn.Select(&transfers, "SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE from_account_id = $1 OR to_account_id = $2 ORDER BY id LIMIT $3 OFFSET $4;", from_account_id, to_account_id, limit, offset)
	if err != nil {
		return nil, err
	}

	return &transfers, nil
}
