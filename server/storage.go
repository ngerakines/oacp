package server

import (
	"context"
	"database/sql"
	"fmt"
	"io"
)

type storage interface {
	RecordLocation(ctx context.Context, state, location string) error
	GetLocation(ctx context.Context, state string) (string, error)

	io.Closer
}

var (
	errNotFound = fmt.Errorf("not found")
	errExists   = fmt.Errorf("exists")
	errInvalid  = fmt.Errorf("invalid")

	errLocationNotFound = fmt.Errorf("location: %w", errNotFound)
	errLocationExists   = fmt.Errorf("location: %w", errExists)
	errLocationInvalid  = fmt.Errorf("location: %w", errInvalid)

	errStateInvalid = fmt.Errorf("state: %w", errInvalid)

	errStorageEngineNotFound = fmt.Errorf("storage engine: %w", errNotFound)
)

type transactionScopedWork func(db *sql.Tx) error

func runTransactionWithOptions(db *sql.DB, txBody transactionScopedWork) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = txBody(tx)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return txErr
		}
		return err
	}
	return tx.Commit()
}