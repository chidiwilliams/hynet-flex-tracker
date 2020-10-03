package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "hynet-flex.db"
	tableName  = "usage"
)

func Connect() (*sql.DB, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return nil, err
	}

	dbDir := homeDir + "/hynet-flex"
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", dbDir+"/"+dbFileName)
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func PrepareDB() (*sql.DB, error) {
	db, err := Connect()
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

func SaveBalance(balance int, db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(balance, date) values(?,?)", tableName))
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(balance, time.Now().UTC())
	return err
}
