package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const file = "test.db"
const create string = `
  CREATE TABLE IF NOT EXISTS infrastructures (
  key VARCHAR NOT NULL PRIMARY KEY,
  url VARCHAR NOT NULL,
  token VARCHAR NOT NULL,
  UNIQUE(url)
  );
  `

func NewDatabase() (*Database, error) {

	// Create connection
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	// initialise table if doesn't exist
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

// database struct so we can implement mutual exclusion
type Database struct {
	db *sql.DB
	mu sync.Mutex
}

// Helper function to avoid repetition
func (data *Database) executeQuery(query string, key string, url string, token string) (sql.Result, error) {

	// Validate parameters
	if key == "" || url == "" || token == "" {
		log.Print("Invalid data entered")
		return nil, nil
	}

	// Execute query
	data.mu.Lock()
	res, err := data.db.Exec(query, key, url, token)
	data.mu.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

// Inserts data into table
func (data *Database) InsertAllCols(key string, url string, token string) (sql.Result, error) {

	// Prepare query
	query := `INSERT INTO infrastructures VALUES (?,?,?);`

	// Execute Query
	res, _ := data.executeQuery(query, key, url, token)

	return res, nil
}

// Deletes row from table. Requires all data for security reasons
func (data *Database) DeleteRow(key string, url string, token string) (sql.Result, error) {

	// Prepare query
	query := `DELETE FROM infrastructures WHERE key = ? AND url = ? AND token = ?;`

	// Execute query
	res, _ := data.executeQuery(query, key, url, token)

	return res, nil
}

// Returns row from table given key
func (data *Database) GetRow(key string) (string, string, error) {

	// Prepare query
	query := `SELECT url, token FROM infrastructures WHERE key = ?;`

	// Execute query
	row := data.db.QueryRow(query, key)

	var url string
	var token string

	if err := row.Scan(&url, &token); err == sql.ErrNoRows {
		log.Printf("Key not found")
		return "", "", err
	}

	return url, token, nil
}

// Return all urls
func (data *Database) GetAllUrls() ([]string, error) {

	// Prepare query
	query := `SELECT url FROM infrastructures;`

	// Execute query
	rows, err := data.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through rows and append data to slice
	var urls []string
	for rows.Next() {
		var i string
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		urls = append(urls, i)
	}
	return urls, nil
}

// Close connection to database
func (data *Database) Close() error {
	err := data.db.Close()
	if err != nil {
		return err
	}
	return nil
}
