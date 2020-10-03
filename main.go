package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/robfig/cron/v3"
)

const (
	defaultFreq = "24h"

	baseURL       = "http://192.168.0.1"
	getProcessURL = baseURL + "/goform/goform_get_cmd_process"
	setProcessURL = baseURL + "/goform/goform_set_cmd_process"
)

var balanceRegex = regexp.MustCompile("\\b([0-9]+)MB\\b")

func main() {
	fmt.Println("MTN HynetFlex Tracker")

	freq, err := promptFrequency()
	if err != nil {
		log.Fatal(err)
	}

	password, err := promptPassword()
	if err != nil {
		log.Fatal(err)
	}

	db, err := prepareDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")

	logIfErr(loginAndSaveBalance(password, db))
	log.Printf("Next run in %s", freq)

	cr := cron.New()
	_, err = cr.AddFunc(fmt.Sprintf("@every %s", freq), func() {
		logIfErr(loginAndSaveBalance(password, db))
		log.Printf("Next run in %s", freq)
	})
	if err != nil {
		log.Fatal(err)
	}

	cr.Start()
	select {}
}

func promptFrequency() (string, error) {
	return (&promptui.Prompt{
		Label:   "Frequency",
		Default: defaultFreq,
		Validate: func(s string) error {
			_, err := time.ParseDuration(s)
			if err != nil {
				return errors.New("invalid duration. See https://golang.org/pkg/time/#ParseDuration for valid values")
			}
			return nil
		},
	}).Run()
}

func promptPassword() (string, error) {
	return (&promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}).Run()
}

func loginAndSaveBalance(adminPassword string, db *sql.DB) error {
	cookies, err := loginToAdmin(adminPassword)
	if err != nil {
		return err
	}

	log.Println("Logged in to admin")

	balance, err := getBalance(cookies)
	if err != nil {
		return err
	}

	log.Printf("Current balance: %dMB\n", balance)

	if err = saveBalance(balance, db); err != nil {
		return err
	}

	log.Println("Balance saved")
	return nil
}

func logIfErr(err error) {
	if err != nil {
		log.Printf("Error: %s", err)
	}
}
