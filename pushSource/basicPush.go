package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	for i := 0; ; i++ {
		// Encode the data
		postBody, _ := json.Marshal(map[string]string{
			"time": time.Now().String(),
		})
		responseBody := bytes.NewBuffer(postBody)

		// Post request to server
		_, err := http.Post("http://127.0.0.1:8080/top/", "application/json", responseBody)
		if err != nil {
			log.Print(err.Error())
		}

		// Print log message and sleep for 5 seconds.
		log.Printf("Sent message %d ", i)
		time.Sleep(5e9)
	}
}
