package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the SQLite database connection
func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	// Limit concurrent connections to prevent database locking
	DB.SetMaxOpenConns(10)
	return DB.Ping()
}

// RunMigrations executes SQL migrations from migrations.sql file
func RunMigrations() error {
	data, err := os.ReadFile("db/migrations.sql")
	if err != nil {
		return err
	}
	_, err = DB.Exec(string(data))
	return err
}
