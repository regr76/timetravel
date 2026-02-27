package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:"+filename+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version INTEGER NOT NULL DEFAULT 1,
		start TEXT NOT NULL,
		end TEXT,
		data TEXT NOT NULL DEFAULT '{}',

		-- Ensure metadata always contains valid JSON
		CHECK (json_valid(data))
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}
