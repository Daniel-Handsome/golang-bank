package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/daniel/master-golang/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandString(3),
		Balance:  int64(utils.RandInt(10, 1000)),
		Currency: utils.RandCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreateAt)

	return account
}

func equalBetween(t *testing.T, a, b Account) {
	require.Equal(t, a.ID, b.ID)
	require.Equal(t, a.Owner, b.Owner)
    require.Equal(t, a.Balance, b.Balance)
    require.Equal(t, a.Currency, b.Currency)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := GetAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	equalBetween(t, account1, account2)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID : account1.ID,
        Balance : int64(utils.RandInt(10, 1000)),
	}

    account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
    require.Equal(t, arg.Balance, account2.Balance)
    require.Equal(t, account1.Currency, account2.Currency)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

    err := testQueries.DeleteAccount(context.Background(), account.ID)

    require.NoError(t, err)

	account, err = testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
    require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}
