package db

import (
	"context"
	"fmt"
	"strings"
)

// create
func (s *Queries) CreateAccount(ctx context.Context, owner string, balance int64, currency string) (*Account, error) {
	row := s.db.QueryRowContext(ctx, "INSERT INTO accounts (owner, balance, currency) VALUES ($1, $2, $3) RETURNING id, owner, balance, currency, created_at;", owner, balance, currency)

	var account Account

	row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)

	return &account, nil
}

// read (id)
func (s *Queries) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	var account Account

	err := s.db.GetContext(ctx, &account, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Queries) GetAccountByIDForUpdate(ctx context.Context, id int64) (*Account, error) {
	var account Account

	err := s.db.GetContext(ctx, &account, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = $1 FOR NO KEY UPDATE;", id)
	if err != nil {
		return nil, err
	}

	return &account, err
}

// read (owner)
func (s *Queries) GetAccountsByOwner(ctx context.Context, owner string) (*[]Account, error) {
	var accounts []Account

	err := s.db.SelectContext(ctx, &accounts, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE owner = $1;", owner)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

// read (owner) (pagination)
func (s *Queries) ListAccounts(ctx context.Context, owner string, limit, offset int64) (*[]Account, error) {
	var accounts []Account

	err := s.db.SelectContext(ctx, &accounts, "SELECT id, owner, balance, currency, created_at FROM accounts WHERE owner = $1 ORDER BY id LIMIT $2 OFFSET $3;", owner, limit, offset)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

// update (for adming use ONLY)
func (s *Queries) UpdateAccount(ctx context.Context, account *Account) (int64, error) {
	return s.db.MustExecContext(ctx, "UPDATE accounts SET owner = $1, balance = $2, currency = $3 WHERE id = $4;", account.Owner, account.Balance, account.Currency, account.ID).RowsAffected()
}

// update owner for accountID
func (s *Queries) UpdateAccountOwner(ctx context.Context, accountID int64, newOwner string) (int64, error) {
	return s.db.MustExecContext(ctx, "UPDATE accounts SET owner = $1 WHERE id = $2;", newOwner, accountID).RowsAffected()
}

// update account balance
func (s *Queries) UpdateAccountBalance(ctx context.Context, id int64, balance int64) (int64, error) {
	return s.db.MustExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2;", balance, id).RowsAffected()
}

// add to account's balance
func (s *Queries) AddAccountBalance(ctx context.Context, id int64, amount int64) (*Account, error) {
	row := s.db.QueryRowContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING id, owner, balance, currency, created_at;", amount, id)

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
func (s *Queries) DeleteAccountByID(ctx context.Context, id int64) (int64, error) {
	return s.db.MustExecContext(ctx, "DELETE FROM accounts WHERE id = $1;", id).RowsAffected()
}
