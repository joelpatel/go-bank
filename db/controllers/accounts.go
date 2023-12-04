package controllers

import (
	"time"

	"github.com/joelpatel/go-bank/db"
)

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   int64     `json:"balance" db:"balance"` // int64 scaled representing last 3 digits as decimal points
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// create
func CreateAccount(owner, balance, currency string) (int64, error) {
	tx := db.DBConn.MustBegin()
	res := tx.MustExec("INSERT INTO accounts (owner, balance, currency) VALUES ($1, $2, $3);", owner, balance, currency)
	err := tx.Commit()

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// read (id)
func GetAccountByID(id string) (*Account, error) {
	var account Account

	err := db.DBConn.Get(&account, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = $1;", id)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

// read (owner)
func GetAccountsByOwner(owner string) (*[]Account, error) {
	var accounts []Account

	err := db.DBConn.Select(&accounts, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE owner = $1;", owner)

	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

// update
func UpdateAccount(account *Account) (int64, error) {
	tx := db.DBConn.MustBegin()
	res := tx.MustExec("UPDATE accounts set owner = $1, balance = $2, currency = $3 WHERE id = $4;", account.Owner, account.Balance, account.Currency, account.ID)
	err := tx.Commit()

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// delete
func DeleteAccountByID(id string) (int64, error) {
	tx := db.DBConn.MustBegin()
	res := tx.MustExec("DELETE FROM accounts WHERE id = $1;", id)
	err := tx.Commit()

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
