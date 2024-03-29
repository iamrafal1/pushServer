package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/iamrafal1/pushServer/db"
)

// Wrapper for handler that generates new infrastructures
func GenerationHandler(dists map[string]*Distributor, data *db.Database) func(http.ResponseWriter, *http.Request) {
	if dists == nil {
		log.Fatal("ERROR nil Distributor session!")
	}
	if data == nil {
		log.Fatal("ERROR nil Database session!")
	}

	// Handler for generating new infrastructures
	return func(w http.ResponseWriter, r *http.Request) {

		for {

			// Generate values
			h1 := hashGenerator()
			key := h1[0:9]
			h2 := hashGenerator()
			url := h2[0:9]
			h3 := hashGenerator()
			token := h3[0:31]

			// Try to insert into db. If fails, generate again
			_, err := data.InsertAllCols(key, url, token)
			if err != nil {
				continue
			}

			// Create distributor
			dist := NewDistributor()
			dists[url] = dist

			// Return information about newly generated information
			fmt.Fprintf(w, `{"key": "%s", "url": "https://127.0.0.1:8080/sub/%s", "token":"%s"}`, key, url, token)
			break
		}

		// Done.
		log.Println("Finished HTTP request at", r.URL.Path)
	}
}

// Wrapper for handling deletions of infrastructure
func DeletionHandler(dists map[string]*Distributor, data *db.Database) func(http.ResponseWriter, *http.Request) {
	if dists == nil {
		log.Fatal("ERROR nil Distributor session!")
	}
	if data == nil {
		log.Fatal("ERROR nil Database session!")
	}

	// Handler for deleting infrastructures
	return func(w http.ResponseWriter, r *http.Request) {

		// Validate key and token
		url := RequestValidator(r, data)
		if url == "" {
			log.Print("Validation failed")
			w.Write([]byte("Validation error"))
			return
		}

		key := r.Header.Get("Push-Key")
		token := r.Header.Get("Push-Token")

		// Remove from database
		_, err := data.DeleteRow(key, url, token)
		if err != nil {
			log.Print("Failed deletion")
			w.Write([]byte("Error"))
			return
		}
		delete(dists, url)

		// Return success message
		w.Write([]byte("Success"))

		// Done.
		log.Println("Finished delete request at", r.URL.Path)
	}
}

// Helper function to generate hashes
func hashGenerator() string {
	// Create random number that is cryptographically safe
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(100000))
	randomString := randomNumber.String()

	// Create timestamp
	timestamp := time.Now().String()

	// Concat timestamp and random string for more security
	fullString := timestamp + randomString

	// Create hash
	h := sha256.New()
	h.Write([]byte(fullString))

	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return string(hashString)
}
