package database

import (
	"os"

	"github.com/jmoiron/sqlx"
)

func ConnectDB() *sqlx.DB {
	AbsPath, err := os.Getwd() // TODO: This
	if err != nil {
		return nil
	}

	DBPath := AbsPath + "/database/database.db"

	connection, err := sqlx.Open("sqlite", DBPath)
	if err != nil {
		return nil
	}

	return connection
}
