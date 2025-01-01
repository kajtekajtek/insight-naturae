// pkg/database/Database.go - database operations
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Init initializes the database
func Init(dbPath, schema string) (*sql.DB, error) {
	// open or create the database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// create the table
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	fmt.Println("Database initialized")
	return db, nil
}