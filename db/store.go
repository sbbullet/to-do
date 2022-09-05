package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}
