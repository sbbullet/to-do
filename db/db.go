package db

import (
	"database/sql"

	"github.com/sbbullet/to-do/util"
)

func NewDB(config *util.Config) *sql.DB {
	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users(
		username TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		full_name TEXT NOT NULL,
		hashed_password TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS tasks(
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		title TEXT NOT NULL,
		is_completed INTEGER DEFAULT 0 CHECK(is_completed IN(0,1)),
		created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (username) REFERENCES users (username) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	return db
}
