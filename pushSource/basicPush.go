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
		req, err := http.NewRequest("POST", "http://127.0.0.1:8080/top", responseBody)
		if err != nil {
			log.Print("Failed to create request")
		}
		req.Header.Add("Push-Key", "1")
		req.Header.Add("Push-Token", "3")

		// Post request to server
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Print(err.Error())
		}

		// Print log message and sleep for 5 seconds.
		log.Printf("Sent message %d ", i)
		time.Sleep(5e9)
	}
}
