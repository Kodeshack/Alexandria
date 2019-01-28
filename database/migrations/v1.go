package migrations

import (
	"database/sql"
	"log"
)

func v1Migration(db *sql.Tx) error {
	sql, err := loadRawSql("1")
	if err != nil {
		log.Fatal("[v1Migration]", err)
		return err
	}

	_, err = db.Exec(sql)

	return err
}
