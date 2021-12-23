package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func Init() {
	var err error
	Database, err = sql.Open("sqlite3", os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}
}
