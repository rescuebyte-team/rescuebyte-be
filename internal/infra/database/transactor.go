package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactor struct {
	db *sqlx.DB
}

func NewTransactor(db *sqlx.DB) Transactor {
	return &transactor{db: db}
}

var _ Transactor = (*transactor)(nil)

func (t *transactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if txFromContext(ctx) != nil {
		return fn(ctx)
	}

	tx, txErr := t.db.BeginTxx(ctx, nil)
	if txErr != nil {
		return fmt.Errorf("begin transaction: %w", txErr)
	}

	defer func() {
		if panicValue := recover(); panicValue != nil {
			_ = tx.Rollback()
			panic(panicValue)
		}
	}()

	txCtx := contextWithTx(ctx, tx)

	if fnErr := fn(txCtx); fnErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(fnErr, fmt.Errorf("rollback transaction: %w", rollbackErr))
		}
		return fnErr
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("commit transaction: %w", commitErr)
	}

	return nil
}

type DBGetter func(ctx context.Context) DB

func NewDBGetter(db *sqlx.DB) DBGetter {
	return func(ctx context.Context) DB {
		if tx := txFromContext(ctx); tx != nil {
			return tx
		}
		return db
	}
}
