// transaction in context of banking system (not database)
package db

import (
	"context"
)

type TransferTxResult struct {
	TransferRecord  Transfer `json:"transfer"`
	FromAccount     Account  `json:"from_account"`
	ToAccount       Account  `json:"to_account"`
	FromEntryRecord Entry    `json:"from_entry"`
	ToEntryRecord   Entry    `json:"to_entry"`
}

// create a transfer record
// create an entry record for: from
// create an entry record for: to
// update balance in account: from
// update balance in account: to
func TransferMoney(ctx context.Context, from_account_id, to_account_id, amount int64) (*TransferTxResult, error) {
	tx := Conn.MustBeginTx(ctx, nil)

	transferRecord, err := CreateTransfer(ctx, tx, from_account_id, to_account_id, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	fromEntry, err := CreateEntry(ctx, tx, from_account_id, -amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	toEntry, err := CreateEntry(ctx, tx, to_account_id, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var fromAccount, toAccount *Account
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
func addAmountInOrder(ctx context.Context, executor Executor, first_account_id, second_account_id, amount int64) (account1, account2 *Account, err error) {
	account1, err = AddAccountBalance(ctx, executor, first_account_id, amount)
	if err != nil {
		return
	}

	account2, err = AddAccountBalance(ctx, executor, second_account_id, -amount)
	return
}
