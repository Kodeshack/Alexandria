package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"tobiwiki.app/database/migrations"
)

func NewDBConnection(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupDB(db *sql.DB) error {
	migrations.LoadAllMigrations()
	return migrations.RunAllMigrations(db)
}
