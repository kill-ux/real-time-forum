// db/database.go
package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database connection and sets up the SQLite database.
// It opens a connection to the specified database file and stores it in the global DB variable.
// Returns an error if the database cannot be opened or initialized.
func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(10)
	return DB.Ping()
}

func RunMigrations() error {
	data, err := os.ReadFile("db/migrations.sql")
	if err != nil {
		return err
	}
	_, err = DB.Exec(string(data))
	return err
}
