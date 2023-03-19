package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const file = "test.db"
const create string = `
  CREATE TABLE IF NOT EXISTS infrastructures (
  key VARCHAR NOT NULL PRIMARY KEY,
  url VARCHAR NOT NULL,
  token VARCHAR NOT NULL,
  UNIQUE(url)
  );`

// Check if basic functioning of db is okay
func DbCheck() error {
	// Open connection to db
	db, err := DbOpen()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get db version to verify SQL works
	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(version)

	// Initialise table if doesn't exist
	err = tableInit()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// Initialise table if doesn't exist
func tableInit() error {
	db, err := DbOpen()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(create); err != nil {
		return err
	}
	return nil
}

// Open connection to database
func DbOpen() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
