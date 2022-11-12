package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joelpatel/go-bank/util"
)

const (
	PRECISION = 6
)

/*
Store provides all functions to execute DB queries and transactions.
*/
type Store struct {
	*Queries // composition ==> to extend functionality of Queries, like inheritance from other lang.
	// all functionalities of Queries will be available to Store.
	db *sql.DB
}

// Contains all input parameters of the transfer tx.
type TransferTxParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

// Result of the transfer tx.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Creates a new Store.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// Executes a function within a database transaction.
func (store *Store) executeTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // using default tx. options
	if err != nil {
		return err
	}

	query := New(tx)
	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr) // combining 2 errors
		}
		return err
	}

	return tx.Commit()
}

/*
Transfer money from one account to another.
- Create a transfer record.
- Create a -ve. account entry for sender.
- Create a +ve. account entry for receiver.
- Subtract values from sender.
- Add values to receiver.
*/
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		// Create a transfer record.
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{ // closure as using result (outer) inside callback func.
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		//Create a -ve. account entry for sender.
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    util.RoundFloat(-arg.Amount, PRECISION),
		})
		if err != nil {
			return err
		}

		// Create a +ve. account entry for receiver.
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    util.RoundFloat(arg.Amount, PRECISION),
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = transferMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 float64,
	accountID2 int64,
	amount2 float64,
) (updatedAccount1 Account, updatedAccount2 Account, err error) {
	account1, err := q.GetAccountForUpdate(ctx, accountID1)
	if err != nil {
		return
	}

	updatedAccount1, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RoundFloat(account1.Balance+amount1, PRECISION),
	})
	if err != nil {
		return
	}

	account2, err := q.GetAccountForUpdate(ctx, accountID2)
	if err != nil {
		return
	}

	updatedAccount2, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      account2.ID,
		Balance: util.RoundFloat(account2.Balance+amount2, PRECISION),
	})
	if err != nil {
		return
	}

	return
}
