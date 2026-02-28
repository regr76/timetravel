package dbutils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableName        = "records"
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		id INTEGER NOT NULL,
		version INTEGER NOT NULL DEFAULT 1,
		start TEXT NOT NULL,
		end TEXT,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY (id, version),

		-- Ensure metadata always contains valid JSON
		CHECK (json_valid(data))
	) STRICT;
	`
)

func InitDB(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:"+filename+"?_journal_mode=WAL&_busy_timeout=1000")
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(1 * time.Second)

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ReadOneVersion(db *sql.DB, id int, version int) (string, error) {
	var idx, ver, start, end, data string
	query := `SELECT * FROM ` + tableName + ` WHERE id = ? AND version = ?`
	err := db.QueryRow(query, id, version).Scan(&idx, &ver, &start, &end, &data)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("{\"id\": %s,\"version\": %s, \"start\": \"%s\", \"end\": \"%s\", \"data\": %s}", idx, ver, start, end, data)
	return result, nil
}

func ReadAllVersions(db *sql.DB, id int) ([]string, error) {
	query := `SELECT * FROM ` + tableName + ` WHERE id = ? ORDER BY version ASC`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var versions []string
	for rows.Next() {
		// scan into all five columns
		var idx, ver, start, end, data string
		if err := rows.Scan(&idx, &ver, &start, &end, &data); err != nil {
			return nil, err
		}
		result := fmt.Sprintf("{\"id\": %s,\"version\": %s, \"start\": \"%s\", \"end\": \"%s\", \"data\": %s}", idx, ver, start, end, data)
		versions = append(versions, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func ReadLatestVersion(db *sql.DB, id int) (string, error) {
	var idx, ver, start, end, data string
	query := `SELECT * FROM ` + tableName + ` WHERE id = ? ORDER BY version DESC LIMIT 1`
	err := db.QueryRow(query, id).Scan(&idx, &ver, &start, &end, &data)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("{\"id\": %s,\"version\": %s, \"start\": \"%s\", \"end\": \"%s\", \"data\": %s}", idx, ver, start, end, data)
	return result, nil
}

func WriteVersion(db *sql.DB, id int, version int, start string, end string, data string) error {
	query := `INSERT INTO ` + tableName + ` (id, version, start, end, data) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, id, version, start, end, data)
	return err
}

func UpdateVersion(db *sql.DB, id int, version int, end string) error {
	query := `UPDATE ` + tableName + ` SET end = ? WHERE id = ? AND version = ?`
	_, err := db.Exec(query, end, id, version)
	return err
}

func ReadAllRows(db *sql.DB) ([]string, error) {
	query := `SELECT * FROM ` + tableName
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var results []string
	for rows.Next() {
		var id, version int
		var start, end, data string
		if err := rows.Scan(&id, &version, &start, &end, &data); err != nil {
			return nil, err
		}
		results = append(results, data)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
