package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iamrafal1/pushServer/db"
	h "github.com/iamrafal1/pushServer/handlers"
)

func main() {

	data, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	distributorMap := initialiseDistributors(data)
	fmt.Println(distributorMap)
	http.HandleFunc("/top", h.WebhookHandler(distributorMap, data))
	http.HandleFunc("/generate", h.GenerationHandler(distributorMap, data))
	http.HandleFunc("/delete", h.DeletionHandler(distributorMap, data))
	http.HandleFunc("/", h.MainHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

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
		http.Handle("/"+url, dist)
		urlMap[url] = dist
	}

	return urlMap
}
