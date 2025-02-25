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
