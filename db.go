package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbPath    = "./hynet-flex.db"
	tableName = "usage"
)

func prepareDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = createTable(db); err != nil {
		return nil, err
	}

	return db, createIndex(db)
}

func createTable(db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    balance integer,
    date datetime
)`, tableName))
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	return err
}

func createIndex(db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS datetime ON %s(date)", tableName))
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	return err
}

func saveBalance(balance int, db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(balance, date) values(?,?)", tableName))
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(balance, time.Now())
	return err
}
