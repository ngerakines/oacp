package server

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
)

type mysqlStorage struct {
	db *sql.DB
}

var _ storage = &mysqlStorage{}

func newMySQLStorage(c *cli.Context) (storage, error) {
	db, err := sql.Open("mysql", c.String("storage-args"))
	if err != nil {
		return nil, err
	}
	return &mysqlStorage{db: db}, nil
}

func (l mysqlStorage) RecordLocation(ctx context.Context, state, location string) error {
	_, err := l.db.ExecContext(ctx, "INSERT INTO state_locations (state, location) VALUES (?, ?)", state, location)
	return err
}

func (l mysqlStorage) GetLocation(ctx context.Context, state string) (string, error) {
	var location string
	txErr := runTransactionWithOptions(l.db, func(tx *sql.Tx) error {
		err := tx.
			QueryRowContext(ctx, "SELECT location FROM state_locations WHERE state = ?", state).
			Scan(&location)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "DELETE FROM locations WHERE state = ?", state)
		return err
	})
	if txErr != nil {
		return "", txErr
	}
	return location, nil
}

func (l mysqlStorage) Close() error {
	return l.db.Close()
}
