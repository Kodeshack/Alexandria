package migrations

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Migration struct {
	Version     int
	Description string
	Migrate     func(*sql.Tx) error
}

func NewMigration(version int, description string, migrateFn func(*sql.Tx) error) *Migration {
	return &Migration{version, description, migrateFn}
}

const schemaVersion = 1

var migrations = []*Migration{
	&Migration{0, "Empty", nil},
}

func addMigration(m *Migration) {
	migrations = append(migrations, m)
}

func loadRawSql(version string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	sql, err := ioutil.ReadFile(filepath.Join(cwd, "database", "migrations", "sql", "migration_v"+version+".sql"))
	if err != nil {
		return "", err
	}

	return string(sql), nil
}

func LoadAllMigrations() {
	addMigration(NewMigration(
		1,
		"Initial Migration",
		v1Migration,
	))
}

// Kindly borrowed from https://github.com/miniflux/miniflux/blob/master/database/migration.go
func RunAllMigrations(db *sql.DB) error {
	var currentVersion int
	db.QueryRow("SELECT VERSION FROM schema_version").Scan(&currentVersion)

	for version := currentVersion + 1; version <= schemaVersion; version++ {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal("[RunAllMigrations]", err)
			return err
		}

		m := migrations[version]

		log.Println("Migrating to version:", m.Version)
		log.Println("Migration", m.Description)

		err = m.Migrate(tx)

		if err != nil {
			tx.Rollback()
			log.Fatal("[RunAllMigrations]", err)
			return err
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			tx.Rollback()
			log.Fatal("[RunAllMigrations]", err)
			return err
		}

		if _, err := tx.Exec(`insert into schema_version (version) values($1)`, version); err != nil {
			tx.Rollback()
			log.Fatal("[Migrate]", err)
			return err
		}

		if err := tx.Commit(); err != nil {
			log.Fatal("[Migrate]", err)
			return err
		}
	}

	return nil
}
