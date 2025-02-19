package test

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB() *sql.DB {
	// Create a temporary database file
	tmpfile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		log.Fatal(err)
	}

	// Open the database
	db, err := sql.Open("sqlite3", tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			german TEXT NOT NULL,
			english TEXT NOT NULL,
			parts TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS words_groups (
			word_id INTEGER,
			group_id INTEGER,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (group_id) REFERENCES groups(id),
			PRIMARY KEY (word_id, group_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
