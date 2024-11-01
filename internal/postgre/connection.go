package postgre

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Connect(databaseDsn *string) error {
	db, err := sql.Open("pgx", *databaseDsn)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func RunInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
