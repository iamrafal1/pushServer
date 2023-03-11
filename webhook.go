package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func webhookHandler(dist *Distributor) func(http.ResponseWriter, *http.Request) {
	if dist == nil {
		panic("nil Distributor session!")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		strData := string(data)
		if err != nil {
			log.Fatal("error reading body")
		}
		if strData != "" {
			dist.messages <- fmt.Sprint(strData)
		}
		// Done.
		log.Println("Finished HTTP request at", r.URL.Path)
	}
}
