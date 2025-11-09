package postgresql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PostgreSQLTransactionAdapter struct {
	db *sqlx.DB
}

func NewPostgreSQLTransactionAdapter(db *sqlx.DB) *PostgreSQLTransactionAdapter {
	return &PostgreSQLTransactionAdapter{db: db}
}

// WithTransaction выполняет функцию в транзакции
func (a *PostgreSQLTransactionAdapter) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	txCtx := context.WithValue(ctx, "tx", tx)
	err = fn(txCtx)

	return err
}
