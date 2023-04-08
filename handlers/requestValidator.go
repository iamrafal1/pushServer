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
	log.Print(key)
	if key == "" {
		return ""
	}
	reqToken := r.Header.Get("Push-Token")
	if reqToken == "" {
		return ""
	}
	log.Print(key, reqToken)

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

// Helper function to deal with CORS. Currently only set to run on cs1
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://cs1.ucc.ie")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Push-Key, Push-Token")
}
