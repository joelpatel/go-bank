package accounts

import (
	"github.com/joelpatel/go-bank/db"
)

// create
func CreateAccount(owner, balance, currency string) (int64, error) {
	tx := db.Conn.MustBegin()
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
func GetAllAccountsPaginated(limit, offset string) (*[]Account, error) {
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
		return 0, err
	}

	return res.RowsAffected()
}
