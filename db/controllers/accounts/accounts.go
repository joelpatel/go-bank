package accounts

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/joelpatel/go-bank/db"
)

const (
	INSERT_ACCOUNT_QUERY = "INSERT INTO accounts (owner, balance, currency) VALUES ($1, $2, $3) RETURNING id, owner, balance, currency, created_at;"
)

func ScanRow(row *sql.Row) (*Account, error) {
	var account Account

	row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)

	return &account, nil
}

// create (self contained transaction)
func CreateAccountTx(owner string, balance int64, currency string) (*Account, error) {
	tx := db.Conn.MustBegin()

	row := tx.QueryRow(INSERT_ACCOUNT_QUERY, owner, balance, currency)

	account, err := ScanRow(row)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return account, nil
}

// create
func CreateAccount(tx *sqlx.Tx, owner string, balance int64, currency string) (*Account, error) {
	return ScanRow(tx.QueryRow(INSERT_ACCOUNT_QUERY, owner, balance, currency))
}

// read (id)
func GetAccountByID(id string) (*Account, error) {
	var account Account

	err := db.Conn.Get(&account, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// read (owner)
func GetAccountsByOwner(owner string) (*[]Account, error) {
	var accounts []Account

	err := db.Conn.Select(&accounts, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE owner = $1;", owner)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

// read all (pagination)
func GetAllAccounts(limit, offset int64) (*[]Account, error) {
	var accounts []Account

	err := db.Conn.Select(&accounts, "SELECT id, owner, balance, currency, created_at FROM accounts ORDER BY id LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

// update
func UpdateAccount(account *Account) (int64, error) {
	tx := db.Conn.MustBegin()
	res := tx.MustExec("UPDATE accounts SET owner = $1, balance = $2, currency = $3 WHERE id = $4;", account.Owner, account.Balance, account.Currency, account.ID)
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return res.RowsAffected()
}

// update account balance
func UpdateAccountBalance(id string, balance int64) (int64, error) {
	tx := db.Conn.MustBegin()
	res := tx.MustExec("UPDATE accounts SET balance = $1 WHERE id = $2;", balance, id)
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return res.RowsAffected()
}

// delete
func DeleteAccountByID(id string) (int64, error) {
	tx := db.Conn.MustBegin()
	res := tx.MustExec("DELETE FROM accounts WHERE id = $1;", id)
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return res.RowsAffected()
}
