// Tests for the account database interface.
package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/joelpatel/go-bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	args := CreateAccontParams{
		Owner:    util.RandomOwner(),
		Balance:  strconv.FormatInt(util.RandomMoney(), 10),
		Currency: util.RandomCurrency(),
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
