package db

import (
	"context"
	"testing"
	"time"

	"github.com/joelpatel/go-bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account *Account) *Entry {
	entryAmount := utils.RandomMoney()
	entry, err := testStore.CreateEntry(context.Background(), account.ID, entryAmount)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, entryAmount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.WithinDuration(t, time.Now(), entry.CreatedAt, 10*time.Second)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t, createRandomAccount(t))
}

func TestGetEntryById(t *testing.T) {
	expectedEntry := createRandomEntry(t, createRandomAccount(t))

	entry, err := testStore.GetEntryByID(context.Background(), expectedEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, expectedEntry.ID, entry.ID)
	require.Equal(t, expectedEntry.AccountID, entry.AccountID)
	require.Equal(t, expectedEntry.Amount, entry.Amount)
	require.Equal(t, expectedEntry.CreatedAt, entry.CreatedAt)
}

func TestGetEntriesByAccountID(t *testing.T) {
	account := createRandomAccount(t)
	expectedEntries := make([]Entry, 10)

	for i := 0; i < 10; i++ {
		expectedEntries[i] = *createRandomEntry(t, account)
	}

	entries, err := testStore.GetEntriesByAccountID(context.Background(), account.ID, 5, 0)
	require.NoError(t, err)

	for i, entry := range *entries {
		require.Equal(t, expectedEntries[i].ID, entry.ID)
		require.Equal(t, expectedEntries[i].AccountID, entry.AccountID)
		require.Equal(t, expectedEntries[i].Amount, entry.Amount)
		require.Equal(t, expectedEntries[i].CreatedAt, entry.CreatedAt)
	}

	entries, err = testStore.GetEntriesByAccountID(context.Background(), account.ID, 5, 5)
	require.NoError(t, err)

	for i, entry := range *entries {
		require.Equal(t, expectedEntries[i+5].ID, entry.ID)
		require.Equal(t, expectedEntries[i+5].AccountID, entry.AccountID)
		require.Equal(t, expectedEntries[i+5].Amount, entry.Amount)
		require.Equal(t, expectedEntries[i+5].CreatedAt, entry.CreatedAt)
	}
}
