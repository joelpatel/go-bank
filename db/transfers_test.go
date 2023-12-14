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

}
