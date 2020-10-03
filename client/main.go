package main

import (
	"log"
	"time"

	db2 "github.com/chidiwilliams/hynet-flex-tracker/db"
)

func main() {
	db, err := db2.Connect()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("SELECT * FROM usage")
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	for result.Next() {
		if err = result.Err(); err != nil {
			log.Fatal(err)
		}

		var balance int
		var date time.Time
		if err = result.Scan(&balance, &date); err != nil {
			log.Fatal(err)
		}

		log.Println(balance, date)
	}
}
