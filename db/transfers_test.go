package db

import (
	"context"
	"testing"
	"time"

	"github.com/joelpatel/go-bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount, toAccount *Account) *Transfer {
	transferAmount := utils.RandomMoney()
	transfer, err := CreateTransfer(context.Background(), testConn, fromAccount.ID, toAccount.ID, transferAmount)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, fromAccount.ID, transfer.FromAccountID)
	require.Equal(t, toAccount.ID, transfer.ToAccountID)
	require.Equal(t, transferAmount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.WithinDuration(t, time.Now(), transfer.CreatedAt, 10*time.Second)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
}

func TestGetTransferByID(t *testing.T) {
	expectedTransfer := createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))

	transfer, err := GetTransferByID(context.Background(), testConn, expectedTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, expectedTransfer.ID, transfer.ID)
	require.Equal(t, expectedTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, expectedTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, expectedTransfer.Amount, transfer.Amount)
	require.Equal(t, expectedTransfer.CreatedAt, transfer.CreatedAt)
}

func TestGetTransfersFrom(t *testing.T) {
	fromAccount := createRandomAccount(t)
	expectedTransfers := make([]Transfer, 10)

	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			expectedTransfers[i/2] = *createRandomTransfer(t, fromAccount, createRandomAccount(t))
		} else {
			createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
		}
	}

	transfers, err := GetTransfersFromTo(context.Background(), testConn, fromAccount.ID, -1, 5, 0)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for i, transfer := range *transfers {
		require.Equal(t, expectedTransfers[i].ID, transfer.ID)
		require.Equal(t, expectedTransfers[i].FromAccountID, transfer.FromAccountID)
		require.Equal(t, expectedTransfers[i].ToAccountID, transfer.ToAccountID)
		require.Equal(t, expectedTransfers[i].Amount, transfer.Amount)
		require.Equal(t, expectedTransfers[i].CreatedAt, transfer.CreatedAt)
	}

	transfers, err = GetTransfersFromTo(context.Background(), testConn, fromAccount.ID, -1, 5, 5)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for i, transfer := range *transfers {
		require.Equal(t, expectedTransfers[i+5].ID, transfer.ID)
		require.Equal(t, expectedTransfers[i+5].FromAccountID, transfer.FromAccountID)
		require.Equal(t, expectedTransfers[i+5].ToAccountID, transfer.ToAccountID)
		require.Equal(t, expectedTransfers[i+5].Amount, transfer.Amount)
		require.Equal(t, expectedTransfers[i+5].CreatedAt, transfer.CreatedAt)
	}
}

func TestGetTransfersTo(t *testing.T) {
	toAccount := createRandomAccount(t)
	expectedTransfers := make([]Transfer, 10)

	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			expectedTransfers[i/2] = *createRandomTransfer(t, createRandomAccount(t), toAccount)
		} else {
			createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
		}
	}

	transfers, err := GetTransfersFromTo(context.Background(), testConn, -1, toAccount.ID, 5, 0)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for i, transfer := range *transfers {
		require.Equal(t, expectedTransfers[i].ID, transfer.ID)
		require.Equal(t, expectedTransfers[i].FromAccountID, transfer.FromAccountID)
		require.Equal(t, expectedTransfers[i].ToAccountID, transfer.ToAccountID)
		require.Equal(t, expectedTransfers[i].Amount, transfer.Amount)
		require.Equal(t, expectedTransfers[i].CreatedAt, transfer.CreatedAt)
	}

	transfers, err = GetTransfersFromTo(context.Background(), testConn, -1, toAccount.ID, 5, 5)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for i, transfer := range *transfers {
		require.Equal(t, expectedTransfers[i+5].ID, transfer.ID)
		require.Equal(t, expectedTransfers[i+5].FromAccountID, transfer.FromAccountID)
		require.Equal(t, expectedTransfers[i+5].ToAccountID, transfer.ToAccountID)
		require.Equal(t, expectedTransfers[i+5].Amount, transfer.Amount)
		require.Equal(t, expectedTransfers[i+5].CreatedAt, transfer.CreatedAt)
	}
}

func TestGetTransfersFromTo(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	var expectedTransfers []Transfer

	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			expectedTransfers = append(expectedTransfers, *createRandomTransfer(t, fromAccount, toAccount))
		} else if i%3 == 0 {
			expectedTransfers = append(expectedTransfers, *createRandomTransfer(t, fromAccount, createRandomAccount(t)))
		} else if i%5 == 0 {
			expectedTransfers = append(expectedTransfers, *createRandomTransfer(t, createRandomAccount(t), toAccount))
		} else {
			createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
		}
	}

	transfers, err := GetTransfersFromTo(context.Background(), testConn, fromAccount.ID, toAccount.ID, 5, 0)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for i, transfer := range *transfers {
		require.Equal(t, expectedTransfers[i].ID, transfer.ID)
		require.Equal(t, expectedTransfers[i].FromAccountID, transfer.FromAccountID)
		require.Equal(t, expectedTransfers[i].ToAccountID, transfer.ToAccountID)
		require.Equal(t, expectedTransfers[i].Amount, transfer.Amount)
		require.Equal(t, expectedTransfers[i].CreatedAt, transfer.CreatedAt)
	}

	transfers, err = GetTransfersFromTo(context.Background(), testConn, fromAccount.ID, toAccount.ID, int64(len(expectedTransfers)-5), 5)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for i, transfer := range *transfers {
		require.Equal(t, expectedTransfers[i+5].ID, transfer.ID)
		require.Equal(t, expectedTransfers[i+5].FromAccountID, transfer.FromAccountID)
		require.Equal(t, expectedTransfers[i+5].ToAccountID, transfer.ToAccountID)
		require.Equal(t, expectedTransfers[i+5].Amount, transfer.Amount)
		require.Equal(t, expectedTransfers[i+5].CreatedAt, transfer.CreatedAt)
	}
}
