package accounts

import (
	"context"
	"fmt"
	"strings"

	"github.com/joelpatel/go-bank/db"
)

// create
func CreateAccount(ctx context.Context, executor db.Executor, owner string, balance int64, currency string) (*Account, error) {
	row := executor.QueryRowContext(ctx, "INSERT INTO accounts (owner, balance, currency) VALUES ($1, $2, $3) RETURNING id, owner, balance, currency, created_at;", owner, balance, currency)

	var account Account

	row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)

	return &account, nil
}

// read (id)
func GetAccountByID(ctx context.Context, executor db.Executor, id string) (*Account, error) {
	var account Account

	err := executor.GetContext(ctx, &account, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func GetAccountByIDForUpdate(ctx context.Context, executor db.Executor, id string) (*Account, error) {
	var account Account

	err := executor.GetContext(ctx, &account, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = $1 FOR NO ID UPDATE;", id)
	if err != nil {
		return nil, err
	}

	return &account, err
}

// read (owner)
func GetAccountsByOwner(ctx context.Context, executor db.Executor, owner string) (*[]Account, error) {
	var accounts []Account

	err := executor.SelectContext(ctx, &accounts, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE owner = $1;", owner)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

// read all (pagination)
func GetAllAccounts(ctx context.Context, executor db.Executor, limit, offset int64) (*[]Account, error) {
	var accounts []Account

	err := executor.SelectContext(ctx, &accounts, "SELECT id, owner, balance, currency, created_at FROM accounts ORDER BY id LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

// update
func UpdateAccount(ctx context.Context, executor db.Executor, account *Account) (int64, error) {
	return executor.MustExecContext(ctx, "UPDATE accounts SET owner = $1, balance = $2, currency = $3 WHERE id = $4;", account.Owner, account.Balance, account.Currency, account.ID).RowsAffected()
}

// update account balance
func UpdateAccountBalance(ctx context.Context, executor db.Executor, id int64, balance int64) (int64, error) {
	return executor.MustExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2;", balance, id).RowsAffected()
}

// add to account's balance
func AddAccountBalance(ctx context.Context, executor db.Executor, id int64, amount int64) (*Account, error) {
	row := executor.QueryRowContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING id, owner, balance, currency, created_at;", amount, id)

	var account Account

	err := row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "balance_nonnegative") {
			// NOTE: may want to get the account to return better formatted string (with actual balance)
			return nil, fmt.Errorf("%d's balance is less than requested amount", id)
		}
		return nil, err
	}

	return &account, nil
}

// delete
func DeleteAccountByID(ctx context.Context, executor db.Executor, id int64) (int64, error) {
	return executor.MustExecContext(ctx, "DELETE FROM accounts WHERE id = $1;", id).RowsAffected()
}
