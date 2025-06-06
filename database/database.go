package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase() {
	db := OpenDatabase()

	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS movies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			release_year INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT NOT NULL,
			ttl INTEGER NOT NULL
		);
	`)

	if err != nil {
		panic(err)
	}
}

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		panic(err)
	}

	return db
}
