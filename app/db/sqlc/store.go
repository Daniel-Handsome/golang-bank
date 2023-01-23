package db

import (
	"context"
	"database/sql"
)

// mock db 所以當server new需要store  去連接資料庫 因此要建立mock db 把store改成interface

type Store interface {
	Querier
	TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
        Queries: New(db),
    }
}