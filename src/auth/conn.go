package main

import (
	"database/sql"
	"log"
	"os"
	"time"
)

const backoff int32 = 2

var connAttempts int32 = 0

func openDb(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDb() *sql.DB {
	dataSource := os.Getenv("DSN")

	for {
		conn, err := openDb(dataSource)
		if err == nil {
			return conn
		}

		log.Println("Waiting for Postgres...")
		connAttempts++

		if connAttempts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(time.Duration(backoff) * time.Second)
		continue
	}
}
