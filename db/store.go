package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// methods to execute CRUD operations in application scope
type Store interface {
	CreateAccount(ctx context.Context, owner string, balance int64, currency string) (*Account, error)
	GetAccountByID(ctx context.Context, id int64) (*Account, error)
	GetAccountByIDForUpdate(ctx context.Context, id int64) (*Account, error)
	GetAccountsByOwner(ctx context.Context, owner string) (*[]Account, error)
	ListAccounts(ctx context.Context, owner string, limit, offset int64) (*[]Account, error)
	UpdateAccount(ctx context.Context, account *Account) (int64, error)
	UpdateAccountBalance(ctx context.Context, id int64, balance int64) (int64, error)
	AddAccountBalance(ctx context.Context, id int64, amount int64) (*Account, error)
	DeleteAccountByID(ctx context.Context, id int64) (int64, error)
	CreateEntry(ctx context.Context, accountID, amount int64) (*Entry, error)
	GetEntryByID(ctx context.Context, id int64) (*Entry, error)
	GetEntriesByAccountID(ctx context.Context, account_id, limit, offset int64) (*[]Entry, error)
	CreateTransfer(ctx context.Context, from_account_id, to_account_id, amount int64) (*Transfer, error)
	GetTransferByID(ctx context.Context, id int64) (*Transfer, error)
	GetTransfersFromTo(ctx context.Context, from_account_id, to_account_id, limit, offset int64) (*[]Transfer, error)
	TransferMoney(ctx context.Context, from_account_id, to_account_id, amount int64) (*TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	conn *sqlx.DB
}

func NewStore(conn *sqlx.DB) Store {
	return &SQLStore{
		Queries: NewQueries(conn),
		conn:    conn,
	}
}
