package db

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joelpatel/go-bank/currency"
	"github.com/joelpatel/go-bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) *Account {
	owner := utils.RandomOwner()
	balance := utils.RandomMoney()

	account, err := testStore.CreateAccount(context.Background(), owner, balance, currency.USD)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, owner, account.Owner)
	require.Equal(t, balance, account.Balance)
	require.Equal(t, currency.USD, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	expectedAccount := createRandomAccount(t)
	account, err := testStore.GetAccountByID(context.Background(), expectedAccount.ID)

	require.NoError(t, err)
	require.Equal(t, expectedAccount.ID, account.ID)
	require.Equal(t, expectedAccount.Owner, account.Owner)
	require.Equal(t, expectedAccount.Balance, account.Balance)
	require.Equal(t, expectedAccount.Currency, account.Currency)
	require.Equal(t, expectedAccount.CreatedAt, account.CreatedAt)
}

func TestGetAccountForUpdate(t *testing.T) {
	expectedAccount := createRandomAccount(t)
	account, err := testStore.GetAccountByIDForUpdate(context.Background(), expectedAccount.ID)

	require.NoError(t, err)
	require.Equal(t, expectedAccount.ID, account.ID)
	require.Equal(t, expectedAccount.Owner, account.Owner)
	require.Equal(t, expectedAccount.Balance, account.Balance)
	require.Equal(t, expectedAccount.Currency, account.Currency)
	require.WithinDuration(t, expectedAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestGetAccountsByOwner(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testStore.CreateAccount(context.Background(), account1.Owner, utils.RandomMoney(), currency.USD)
	require.NoError(t, err)

	accounts, err := testStore.GetAccountsByOwner(context.Background(), account1.Owner)
	require.NoError(t, err)

	require.Equal(t, 2, len(*accounts))
	require.Contains(t, []int64{(*accounts)[0].ID, (*accounts)[1].ID}, account1.ID)
	require.Contains(t, []int64{(*accounts)[0].ID, (*accounts)[1].ID}, account2.ID)
	require.Equal(t, account1.Owner, (*accounts)[0].Owner)
	require.Equal(t, account1.Owner, (*accounts)[1].Owner)
	require.Contains(t, []int64{(*accounts)[0].Balance, (*accounts)[1].Balance}, account1.Balance)
	require.Contains(t, []int64{(*accounts)[0].Balance, (*accounts)[1].Balance}, account2.Balance)
	require.Equal(t, account1.Currency, (*accounts)[0].Currency)
	require.Equal(t, account2.Currency, (*accounts)[1].Currency)
	require.Contains(t, []time.Time{(*accounts)[0].CreatedAt, (*accounts)[1].CreatedAt}, account1.CreatedAt)
	require.Contains(t, []time.Time{(*accounts)[0].CreatedAt, (*accounts)[1].CreatedAt}, account2.CreatedAt)
}

func TestGetAllAccounts(t *testing.T) {
	expectedAccounts := make([]Account, 10)

	expectedAccounts[0] = *createRandomAccount(t)

	for i := 1; i < 10; i++ {
		account, err := testStore.CreateAccount(context.Background(), expectedAccounts[0].Owner, utils.RandomMoney(), currency.USD)
		expectedAccounts[i] = *account
		require.NoError(t, err)
	}

	accounts, err := testStore.ListAccounts(context.Background(), expectedAccounts[0].Owner, 5, 0)

	require.NoError(t, err)
	for i := 0; i < 5; i++ {
		require.Equal(t, expectedAccounts[i].ID, (*accounts)[i].ID)
		require.Equal(t, expectedAccounts[i].Owner, (*accounts)[i].Owner)
		require.Equal(t, expectedAccounts[i].Balance, (*accounts)[i].Balance)
		require.Equal(t, expectedAccounts[i].Currency, (*accounts)[i].Currency)
		require.WithinDuration(t, expectedAccounts[i].CreatedAt, (*accounts)[i].CreatedAt, time.Second)
	}

	accounts, err = testStore.ListAccounts(context.Background(), expectedAccounts[0].Owner, 5, 5)
	require.NoError(t, err)
	for i, j := 0, 5; i < 5 && j < 10; i, j = i+1, j+1 {
		require.Equal(t, expectedAccounts[j].ID, (*accounts)[i].ID)
		require.Equal(t, expectedAccounts[j].Owner, (*accounts)[i].Owner)
		require.Equal(t, expectedAccounts[j].Balance, (*accounts)[i].Balance)
		require.Equal(t, expectedAccounts[j].Currency, (*accounts)[i].Currency)
		require.WithinDuration(t, expectedAccounts[j].CreatedAt, (*accounts)[i].CreatedAt, time.Second)
	}
}

func TestUpdateAccount(t *testing.T) {
	originalAccount := createRandomAccount(t)

	var updatedCurrency string
	if originalAccount.Currency == currency.USD {
		updatedCurrency = currency.INR
	} else {
		updatedCurrency = currency.USD
	}

	expectedAccount := Account{
		ID:        originalAccount.ID,
		Owner:     "new_owner",
		Balance:   2000, // 2000 not possible via random amount generator
		Currency:  updatedCurrency,
		CreatedAt: originalAccount.CreatedAt,
	}

	rowsAffected, err := testStore.UpdateAccount(context.Background(), &expectedAccount)
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)

	updatedAccount, err := testStore.GetAccountByID(context.Background(), originalAccount.ID)
	require.NoError(t, err)
	require.Equal(t, expectedAccount.ID, updatedAccount.ID)
	require.Equal(t, expectedAccount.Owner, updatedAccount.Owner)
	require.Equal(t, expectedAccount.Balance, updatedAccount.Balance)
	require.Equal(t, expectedAccount.Currency, updatedAccount.Currency)
	require.Equal(t, expectedAccount.CreatedAt, updatedAccount.CreatedAt)
}

func TestUpdateAccountOwner(t *testing.T) {
	account := createRandomAccount(t)

	rowsAffected, err := testStore.UpdateAccountOwner(context.Background(), account.ID, "new_owner")
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)

	updatedAccout, err := testStore.GetAccountByID(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.ID, updatedAccout.ID)
	require.Equal(t, "new_owner", updatedAccout.Owner)
	require.Equal(t, account.Balance, updatedAccout.Balance)
	require.Equal(t, account.Currency, updatedAccout.Currency)
	require.Equal(t, account.CreatedAt, updatedAccout.CreatedAt)
}

func TestUpdateAccountBalance(t *testing.T) {
	originalAccount := createRandomAccount(t)

	rowsAffected, err := testStore.UpdateAccountBalance(context.Background(), originalAccount.ID, originalAccount.Balance+2000)
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)

	updatedAccount, err := testStore.GetAccountByID(context.Background(), originalAccount.ID)
	require.NoError(t, err)
	require.Equal(t, originalAccount.ID, updatedAccount.ID)
	require.Equal(t, originalAccount.Owner, updatedAccount.Owner)
	require.Equal(t, originalAccount.Balance+2000, updatedAccount.Balance)
	require.Equal(t, originalAccount.Currency, updatedAccount.Currency)
	require.Equal(t, originalAccount.CreatedAt, updatedAccount.CreatedAt)
}

func TestAddAccountBalance(t *testing.T) {
	originalAccount := createRandomAccount(t)

	updatedAccount, err := testStore.AddAccountBalance(context.Background(), originalAccount.ID, 2000)

	require.NoError(t, err)
	require.Equal(t, originalAccount.ID, updatedAccount.ID)
	require.Equal(t, originalAccount.Owner, updatedAccount.Owner)
	require.Equal(t, originalAccount.Balance+2000, updatedAccount.Balance)
	require.Equal(t, originalAccount.Currency, updatedAccount.Currency)
	require.Equal(t, originalAccount.CreatedAt, updatedAccount.CreatedAt)

	_, err = testStore.AddAccountBalance(context.Background(), originalAccount.ID, -10000) // random generate max 1000 + 2000 leads to max 3000 ==> this should lead to negative amount
	expectedError := fmt.Errorf("%d's balance is less than requested amount", originalAccount.ID)
	require.Error(t, expectedError, err)
}

func TestDeleteAccountByID(t *testing.T) {
	originalAccount := createRandomAccount(t)

	rowsAffected, err := testStore.DeleteAccountByID(context.Background(), originalAccount.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)

	account, err := testStore.GetAccountByID(context.Background(), originalAccount.ID)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), pgx.ErrNoRows.Error()))
	require.Empty(t, account)
}
