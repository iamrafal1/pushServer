package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/iamrafal1/pushServer/db"
)

// Wrapper for handler function that requires a database and a distributor
func WebhookHandler(dists map[string]*Distributor, data *db.Database) func(http.ResponseWriter, *http.Request) {
	if dists == nil {
		log.Fatal("ERROR nil Distributor session!")
	}
	if data == nil {
		log.Fatal("ERROR nil Database session!")
	}

	// Handler for incoming webhooks
	return func(w http.ResponseWriter, r *http.Request) {

		// Validate key and token
		url := requestValidator(r, data)
		if url == "" {
			log.Print("Validation failed")
			return
		}

		// Read message body
		message, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("error reading body")
		}
		strData := string(message)
		if strData == "" {
			log.Fatal("body empty")
		}
		dists[url].messages <- fmt.Sprint(strData)

		// Done.
		log.Println("Finished HTTP request at", r.URL.Path)
	}
}

// Helper function for validating a request and determining correct distributor
func requestValidator(r *http.Request, d *db.Database) string {

	// Get info from request header
	key := r.Header.Get("Push-Key")
	if key == "" {
		return ""
	}
	reqToken := r.Header.Get("Push-Token")
	if reqToken == "" {
		return ""
	}

	// retrieve data from db and compare with header info
	url, actualToken, err := d.GetRow(key)
	if err != nil {
		log.Print("Db Error")
		return ""
	}
	if actualToken != reqToken {
		log.Print("Authentication Failed")
		return ""
	}

	return url

}
