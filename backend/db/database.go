// db/database.go
package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(10)
	return nil
}

func RunMigrations() error {
	data, _ := os.ReadFile("db/migrations.sql")
	_, err := DB.Exec(string(data))
	return err
}
