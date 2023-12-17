package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferMoney(context.Background(), account1.ID, account2.ID, amount)
			errs <- err
			results <- *result
		}()
	}

	// check results
	ith := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.TransferRecord
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.WithinDuration(t, time.Now(), transfer.CreatedAt, 10*time.Second)
		require.NotZero(t, transfer.ID)
		_, err = testStore.GetTransferByID(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntryRecord
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.WithinDuration(t, time.Now(), fromEntry.CreatedAt, 10*time.Second)
		_, err = testStore.GetEntryByID(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntryRecord
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.WithinDuration(t, time.Now(), toEntry.CreatedAt, 10*time.Second)
		_, err = testStore.GetEntryByID(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.Equal(t, account1.Owner, fromAccount.Owner)
		require.Equal(t, account1.Currency, fromAccount.Currency)
		require.Equal(t, account1.CreatedAt, fromAccount.CreatedAt)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		require.Equal(t, account2.Owner, toAccount.Owner)
		require.Equal(t, account2.Currency, toAccount.Currency)
		require.Equal(t, account2.CreatedAt, toAccount.CreatedAt)

		// check balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.Positive(t, diff1)
		require.Zero(t, diff1%amount) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, ith, k)
		ith[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := testStore.TransferMoney(context.Background(), fromAccountID, toAccountID, amount)
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the updated final balance
	updatedAccount1, err := testStore.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
