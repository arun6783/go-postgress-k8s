package db

import (
	"context"
	"database/sql"
	"fmt"
)

//store provides all function to execute db in transactions
//this is a example of composition instead of using inheritance
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//takes in context and a callback fun as parameteres
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//tranfer performs a money transfer from one account to other
// it creates a transfer record, add account entries, updates account blance with in a single database transaction

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.Createtransfer(ctx, CreatetransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		var account1ID, account1Amount, account2Id, account2Amount int64

		if arg.FromAccountID < arg.ToAccountID {
			account1ID = arg.FromAccountID
			account1Amount = -arg.Amount
			account2Id = arg.ToAccountID
			account2Amount = arg.Amount

		} else {
			account1ID = arg.ToAccountID
			account1Amount = arg.Amount
			account2Id = arg.FromAccountID
			account2Amount = -arg.Amount
		}

		transferResult, transferErr := transferAmount(ctx, q, account1ID, account1Amount, account2Id, account2Amount)

		if transferErr != nil {
			return err
		}
		result.FromAccount = transferResult.FromAccount
		result.ToAccount = transferResult.ToAccount

		return nil
	})
	return result, err
}

func transferAmount(ctx context.Context, q *Queries, fromAccountId int64, fromAccountAmount int64, toAccountId int64, toAccountAmount int64) (TransferTxResult, error) {
	var result TransferTxResult
	var err error

	result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountId,
		Amount: fromAccountAmount,
	})

	if err != nil {
		return result, err
	}
	result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountId,
		Amount: toAccountAmount,
	})

	if err != nil {
		return result, err
	}

	return result, nil
}
