package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/arun6783/go-postgress-k8s/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {

	account := createRandomAccount(t)

	accountFromDb, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountFromDb)
	require.Equal(t, account.Owner, accountFromDb.Owner)
	require.Equal(t, account.Balance, accountFromDb.Balance)
	require.Equal(t, account.Currency, accountFromDb.Currency)
	require.WithinDuration(t, account.CreatedAt, accountFromDb.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomMoney(),
	}
	accountFromDb, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, accountFromDb)
	require.Equal(t, account.Owner, accountFromDb.Owner)
	//here testing updated result
	require.Equal(t, arg.Balance, accountFromDb.Balance)

	require.Equal(t, account.Currency, accountFromDb.Currency)
	require.WithinDuration(t, account.CreatedAt, accountFromDb.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {

	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	accountFromDb, getAccountErr := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, getAccountErr)
	require.EqualError(t, getAccountErr, sql.ErrNoRows.Error())
	require.Empty(t, accountFromDb)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
