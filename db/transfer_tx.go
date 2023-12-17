// transaction in context of banking system (not database)
package db

import (
	"context"
)

// create a transfer record
// create an entry record for: from
// create an entry record for: to
// update balance in account: from
// update balance in account: to
func (s *SQLStore) TransferMoney(ctx context.Context, from_account_id, to_account_id, amount int64) (*TransferTxResult, error) {
	tx := s.conn.MustBeginTx(ctx, nil)

	q := NewQueries(tx)

	transferRecord, err := q.CreateTransfer(ctx, from_account_id, to_account_id, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	fromEntry, err := q.CreateEntry(ctx, from_account_id, -amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	toEntry, err := q.CreateEntry(ctx, to_account_id, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var fromAccount, toAccount *Account
	if from_account_id < to_account_id {
		fromAccount, toAccount, err = addAmountInOrder(ctx, q, from_account_id, to_account_id, -amount)
	} else {
		fromAccount, toAccount, err = addAmountInOrder(ctx, q, to_account_id, from_account_id, amount)
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
func addAmountInOrder(ctx context.Context, q *Queries, first_account_id, second_account_id, amount int64) (account1, account2 *Account, err error) {
	account1, err = q.AddAccountBalance(ctx, first_account_id, amount)
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, second_account_id, -amount)
	return
}
