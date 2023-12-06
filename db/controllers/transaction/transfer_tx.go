// transaction in context of banking system (not database)
package transaction

import (
	"context"

	"github.com/joelpatel/go-bank/db"
	"github.com/joelpatel/go-bank/db/controllers/accounts"
	"github.com/joelpatel/go-bank/db/controllers/entries"
	"github.com/joelpatel/go-bank/db/controllers/transfers"
)

type TransferTxResult struct {
	TransferRecord  transfers.Transfer `json:"transfer"`
	FromAccount     accounts.Account   `json:"from_account"`
	ToAccount       accounts.Account   `json:"to_account"`
	FromEntryRecord entries.Entry      `json:"from_entry"`
	ToEntryRecord   entries.Entry      `json:"to_entry"`
}

// create a transfer record
// create an entry record for: from
// create an entry record for: to
// update balance in account: from
// update balance in account: to
func TransferMoney(ctx context.Context, from_account_id, to_account_id, amount int64) (*TransferTxResult, error) {
	tx := db.Conn.MustBeginTx(ctx, nil)

	transferRecord, err := transfers.CreateTransfer(ctx, tx, from_account_id, to_account_id, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	fromEntry, err := entries.CreateEntry(ctx, tx, from_account_id, -amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	toEntry, err := entries.CreateEntry(ctx, tx, to_account_id, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var fromAccount, toAccount *accounts.Account
	if from_account_id < to_account_id {
		fromAccount, toAccount, err = addAmountInOrder(ctx, tx, from_account_id, to_account_id, -amount)
	} else {
		fromAccount, toAccount, err = addAmountInOrder(ctx, tx, to_account_id, from_account_id, amount)
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &TransferTxResult{
		TransferRecord:  *transferRecord,
		FromEntryRecord: *fromEntry,
		ToEntryRecord:   *toEntry,
		FromAccount:     *fromAccount,
		ToAccount:       *toAccount,
	}, nil
}

// add +amount to first_account
// add -amount to second_account
func addAmountInOrder(ctx context.Context, executor db.Executor, first_account_id, second_account_id, amount int64) (account1, account2 *accounts.Account, err error) {
	account1, err = accounts.AddAccountBalance(ctx, executor, first_account_id, amount)
	if err != nil {
		return
	}

	account2, err = accounts.AddAccountBalance(ctx, executor, second_account_id, -amount)
	return
}
