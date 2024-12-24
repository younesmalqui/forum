package config

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() error {
	db, err := sql.Open(DRIVER_NAME, "./database/"+DATABASE_NAME)
	if err != nil {
		return NewInternalError(err)
	}
	DB = db
	return nil
}
