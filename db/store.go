package db

import (
	"database/sql"
	"log"
)

type Store struct {
	DB *sql.DB
}

func NewStore() *Store {
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal("cannot open database. Error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("cannot find the running database. Error: ", err)
	}

	store := &Store{
		DB: db,
	}

	if err := store.migrateDatabase(); err != nil {
		log.Fatal("cannot migrate database. Error: ", err)
	}

	return store
}

func (store *Store) migrateDatabase() error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users(
		username TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		full_name TEXT NOT NULL,
		hashed_passworf TEXT NOT NULL,
		created_at INTEGER NOT NULL DEFAULT julianday('now'),
	);
	CREATE TABLE IF NOT EXISTS tasks(
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		title TEXT NOT NULL,
		date INTEGER NOT NULL,
		start_time INTEGER NOT NULL,
		end_time INTEGER NOT NULL,
		is_completed INTEGER DEFAULT 0,
    FOREIGN KEY (username) REFERENCES users (username) ON DELETE CASCADE
	);
	`

	_, err := store.DB.Exec(sqlStmt)
	return err
}
