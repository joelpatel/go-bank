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

	fromAccount, err := accounts.AddAccountBalance(ctx, tx, from_account_id, -amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	toAccount, err := accounts.AddAccountBalance(ctx, tx, to_account_id, amount)
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
