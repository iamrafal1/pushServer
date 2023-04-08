package main

import (
	"log"
	"net/http"

	"github.com/iamrafal1/pushServer/db"
	h "github.com/iamrafal1/pushServer/handlers"
)

func main() {

	// Open database
	data, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	distributorMap := initialiseDistributors(data)

	// Initialise routes
	http.HandleFunc("/top", h.WebhookHandler(distributorMap, data))
	http.HandleFunc("/generate", h.GenerationHandler(distributorMap, data))
	http.HandleFunc("/delete", h.DeletionHandler(distributorMap, data))
	http.HandleFunc("/index", h.MainHandler)
	http.HandleFunc("/sub/", h.SubHandler(distributorMap))

	// Assuming there is a server.crt and server.key file existing in the local directory, run TLS server
	log.Fatal(http.ListenAndServeTLS("localhost:8080", "server.crt", "server.key", nil))
}

// Helper function to read and create distributors from database
func initialiseDistributors(d *db.Database) map[string]*h.Distributor {

	// Get urls from database
	urls, err := d.GetAllUrls()
	if err != nil {
		log.Fatal(err)
	}

	// Create url - dist map
	urlMap := make(map[string]*h.Distributor)

	// Iterate through distributors and initialise them
	for _, url := range urls {
		dist := h.NewDistributor()
		log.Print(url)
		urlMap[url] = dist
	}

	return urlMap
}
