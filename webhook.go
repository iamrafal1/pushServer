package main

import (
	"io"
	"log"
	"net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {

	_, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading body")
	}

	// need to encapsulate handler in another function so we can pass extra data to it

	// Done.
	log.Println("Finished HTTP request at", r.URL.Path)
}
