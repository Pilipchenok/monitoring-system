package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
