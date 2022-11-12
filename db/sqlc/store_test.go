package db

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/joelpatel/go-bank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">> before:\t%v\t%v\n", account1.Balance, account2.Balance)

	n := 5
	amount := float64(10.00001)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	// run n go routines to exhaustively test database transactions are 100% working correctly
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errors <- err
			results <- result
		}()
	}

	// check result and err
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)        // auto increment field => thus should not be zero
		require.NotZero(t, transfer.CreatedAt) // timestamp of when it was executed (by db)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check account entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check account entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check subtracted values from sender
		fmt.Printf(">> tx:\t%v\t%v\n", fromAccount.Balance, toAccount.Balance)
		diff1 := util.RoundFloat(account1.Balance-fromAccount.Balance, 6)
		// fmt.Printf(">> %v - %v = %v\n", account1.Balance, fromAccount.Balance, diff1)
		require.NoError(t, err)
		diff2 := util.RoundFloat(toAccount.Balance-account2.Balance, 6)
		require.NoError(t, err)
		require.Equal(t, diff1, diff2)
		require.NoError(t, err)
		require.True(t, diff1 > 0)

		// fmt.Printf(">> diff = %v\tamount = %v\n", diff1, amount)
		// fmt.Printf("diff1 %% amount = %v\n", util.RoundFloat(math.Mod(diff1, amount), PRECISION))
		require.True(t, util.RoundFloat(math.Mod(diff1, amount), PRECISION) == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
		// check added values to receiver
	}

	// check final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">> after:\t%v\t%v\n", updatedAccount1.Balance, updatedAccount2.Balance)

	// account1.Balance - n * amountFl == updatedAccount1.Balance
	require.NoError(t, err)
	require.Equal(t, account1.Balance-float64(n)*amount, updatedAccount1.Balance)
	// account2.Balance + n * amountFl == updatedAccount2.Balance
	require.NoError(t, err)
	require.Equal(t, account2.Balance+float64(n)*amount, updatedAccount2.Balance)
}
