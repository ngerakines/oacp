package server

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

type postgresStorage struct {
	db *sql.DB
}

var _ storage = &postgresStorage{}

func newPGStorage(c *cli.Context) (storage, error) {
	db, err := sql.Open("postgres", c.String("storage-args"))
	if err != nil {
		return nil, err
	}
	return &postgresStorage{db: db}, nil
}

func (l postgresStorage) RecordLocation(ctx context.Context, state, location string) error {
	_, err := l.db.ExecContext(ctx, "INSERT INTO state_locations (state, location) VALUES ($1, $2)", state, location)
	return err
}

func (l postgresStorage) GetLocation(ctx context.Context, state string) (string, error) {
	var location string
	txErr := runTransactionWithOptions(l.db, func(tx *sql.Tx) error {
		err := tx.
			QueryRowContext(ctx, "SELECT location FROM state_locations WHERE state = $1", state).
			Scan(&location)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "DELETE FROM locations WHERE state = $1", state)
		return err
	})
	if txErr != nil {
		return "", txErr
	}
	return location, nil
}

func (l postgresStorage) Close() error {
	return l.db.Close()
}
