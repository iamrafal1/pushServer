package handlers

import (
	"log"
	"net/http"

	"github.com/iamrafal1/pushServer/db"
)

// Helper function for validating a request and determining correct distributor
func RequestValidator(r *http.Request, d *db.Database) string {

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
