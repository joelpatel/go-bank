// Tests for the account database interface.
package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	args := CreateAccontParams{
		Owner:    "John Doe",
		Balance:  "1000000",
		Currency: "USD",
	}

	account, err := testQueries.CreateAccont(context.Background(), args)
	require.NoError(t, err) // fails the test if err != nil
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
