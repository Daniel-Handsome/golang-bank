package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
        Queries: New(db),
    }
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := New(tx)
	
	if err = fn(query); err != nil {
		// 這邊要重新給一個 不然會return返回錯了error message
        if rbErr := tx.Rollback(); nil != rbErr {
			return fmt.Errorf("tx error: %v", err)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FormAccount Account `json:"form_account"`
	ToAccount Account `json:"to_account"`
	FormEntry Entry `json:"form_entry"`
	ToEntry Entry `json:"to_entry"`
}


func (store *Store) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.Createtransfer(ctx, CreatetransferParams{
			FromAccountID: params.FromAccountID,
            ToAccountID:   params.ToAccountID,
            Amount:        params.Amount,
		})
		if err != nil {
			return err
		}

		result.FormEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount: -params.Amount,
		})
		if err!= nil {
            return err
        }

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
            Amount: params.Amount,
		})

		if err != nil {
			return err
		}

		// get account and update
		// fromAccount, err := q.GetAccountForUpdate(ctx, params.FromAccountID)
		// if err != nil {
		// 	return err
		// }
	
		

		// get account and update
		// toAccount, err := q.GetAccountForUpdate(ctx, params.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// 防止帳戶戶轉 預防race 最好都是取同一個鎖
		if params.FromAccountID > params.ToAccountID {
			result.FormAccount, result.ToAccount, err = addAmount(
				ctx, q, params.FromAccountID, -params.Amount, params.ToAccountID, params.Amount,
			)
		}else {
			result.FormAccount, result.ToAccount, err = addAmount(
				ctx, q, params.ToAccountID, params.Amount, params.FromAccountID, -params.Amount,
			)
		}

		return err
		
		// result.FormAccount, err = q.UpdateAccountByBalance(ctx, UpdateAccountByBalanceParams{
		// 	ID: params.FromAccountID,
		// 	Amount:  - params.Amount,
		// })

		// result.ToAccount, err = q.UpdateAccountByBalance(ctx, UpdateAccountByBalanceParams{
		// 	ID: params.ToAccountID,
		// 	Amount:  params.Amount,
		// })
	})


	return result, err
}

func addAmount(
	ctx context.Context,
	q *Queries,
	firstAccuntId int64,
	firstAmount int64,
	secondAccountId int64,
	secondAmount int64,
) (firstAccunt Account, secondAccount Account, err error) {
	firstAccunt, err =q.UpdateAccountByBalance(ctx, UpdateAccountByBalanceParams{
		ID: firstAccuntId,
		Amount: firstAmount,
	})
	
	if err!= nil {
        return
    }

	secondAccount, err =q.UpdateAccountByBalance(ctx, UpdateAccountByBalanceParams{
		ID: secondAccountId,
		Amount: secondAmount,
	})
	
	return
}