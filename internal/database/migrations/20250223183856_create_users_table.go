package database

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up00001, Down00001)
}

func Up00001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "CREATE TABLE users (id VARCHAR(255) PRIMARY KEY, login VARCHAR(255) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL, registerDate TIMESTAMP NOT NULL);")
	return err
}

func Down00001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE users;")
	return err
}
