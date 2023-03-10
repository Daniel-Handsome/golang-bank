// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: transfer.sql

package db

import (
	"context"
)

const createtransfer = `-- name: Createtransfer :one
INSERT INTO transfer (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING id, from_account_id, to_account_id, amount, create_at
`

type CreatetransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (q *Queries) Createtransfer(ctx context.Context, arg CreatetransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createtransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreateAt,
	)
	return i, err
}

const deletetransfer = `-- name: Deletetransfer :exec
DELETE FROM transfer WHERE id = $1
`

func (q *Queries) Deletetransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletetransfer, id)
	return err
}

const gettransfer = `-- name: Gettransfer :one
SELECT id, from_account_id, to_account_id, amount, create_at FROM transfer
WHERE id = $1 LIMIT 1
`

func (q *Queries) Gettransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, gettransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreateAt,
	)
	return i, err
}

const gettransfers = `-- name: Gettransfers :many
SELECT id, from_account_id, to_account_id, amount, create_at FROM transfer
ORDER BY id
LIMIT $1
OFFSET $2
`

type GettransfersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) Gettransfers(ctx context.Context, arg GettransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, gettransfers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreateAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatetransfer = `-- name: Updatetransfer :one
UPDATE transfer
SET amount = $2
WHERE id = $1
RETURNING id, from_account_id, to_account_id, amount, create_at
`

type UpdatetransferParams struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

func (q *Queries) Updatetransfer(ctx context.Context, arg UpdatetransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, updatetransfer, arg.ID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreateAt,
	)
	return i, err
}
