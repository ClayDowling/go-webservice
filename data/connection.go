package data

import (
	"database/sql"
	"log"
)

const DATABASE_DRIVER = "sqlite"
const DATABASE_CONNECTION = "users.db"

func Connect() (*sql.DB, error) {
	db, err := sql.Open(DATABASE_DRIVER, DATABASE_CONNECTION)
	if err != nil {
		log.Printf("Opening database: %s", err)
	}
	return db, err
}
