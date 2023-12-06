package entries

import (
	"context"

	"github.com/joelpatel/go-bank/db"
)

// create
func CreateEntry(ctx context.Context, executor db.Executor, accountID, amount int64) (*Entry, error) {
	row := executor.QueryRowContext(ctx, "INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING id, account_id, amount, created_at;", accountID, amount)

	var entry Entry

	err := row.Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// read
func GetEntryByID(ctx context.Context, executor db.Executor, id string) (*Entry, error) {
	var entry Entry

	err := executor.GetContext(ctx, &entry, "SELECT id, account_id, amount, created_at FROM entries WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// read all for account_id (pagination)
func GetEntriesByAccountID(ctx context.Context, executor db.Executor, account_id, limit, offset int64) (*[]Entry, error) {
	var entries []Entry

	err := executor.SelectContext(ctx, &entries, "SELECT id, account_id, amount, created_at FROM entries WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3;", account_id, limit, offset)
	if err != nil {
		return nil, err
	}

	return &entries, nil
}
