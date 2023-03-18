package main

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"time"
)

func hashGenerator() string {
	// Create random number that is cryptographically safe
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(100000))
	randomString := randomNumber.String()

	// Create timestamp
	timestamp := time.Now().String()

	// Concat timestamp and random string for more security
	fullString := timestamp + randomString

	// Create hash
	h := sha256.New()
	h.Write([]byte(fullString))
	hashString := h.Sum(nil)
	return string(hashString)
}
