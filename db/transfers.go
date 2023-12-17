package db

import (
	"context"
)

// create
func (s *Queries) CreateTransfer(ctx context.Context, from_account_id, to_account_id, amount int64) (*Transfer, error) {
	row := s.db.QueryRowContext(ctx, "INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3) RETURNING id, from_account_id, to_account_id, amount, created_at;", from_account_id, to_account_id, amount)

	var transfer Transfer

	err := row.Scan(&transfer.ID, &transfer.FromAccountID, &transfer.ToAccountID, &transfer.Amount, &transfer.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// read (id)
func (s *Queries) GetTransferByID(ctx context.Context, id int64) (*Transfer, error) {
	var transfer Transfer

	err := s.db.GetContext(ctx, &transfer, "SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// read (from_account_id OR to_account_id)
// (-1 if don't want to search for from exor to)
// paginated
func (s *Queries) GetTransfersFromTo(ctx context.Context, from_account_id, to_account_id, limit, offset int64) (*[]Transfer, error) {
	var transfers []Transfer

	err := s.db.SelectContext(ctx, &transfers, "SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE from_account_id = $1 OR to_account_id = $2 ORDER BY id LIMIT $3 OFFSET $4;", from_account_id, to_account_id, limit, offset)
	if err != nil {
		return nil, err
	}

	return &transfers, nil
}
