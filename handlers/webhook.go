package handlers

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
		EnableCors(&w)
		// Validate key and token
		url := RequestValidator(r, data)
		if url == "" {
			log.Print("Validation failed")
			w.Write([]byte("Validation failed"))
			return
		}

		// Read message body
		message, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print("error reading body")
			w.Write([]byte("Validation failed"))
			return
		}
		strData := string(message)
		log.Print(strData)
		if strData == "" {
			log.Print("body empty")
			w.Write([]byte("No message entered"))
			return
		}

		// Write to the appropriate distributors message channel
		dists[url].messages <- fmt.Sprint(strData)

		// Done.
		log.Println("Finished HTTP request at", r.URL.Path)
		w.Write([]byte("Success"))
	}
}
