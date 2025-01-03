// pkg/database/dbutils.go - generic database operations
package dbutils

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Init initializes the database
func Init(dbPath string) (*sql.DB, error) {
	// open or create the database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	fmt.Println("Database initialized")
	return db, nil
}