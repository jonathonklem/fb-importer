package main

import (
	"log"
	"os"
)

func getToken() string {
	token, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("Error reading token")
		os.Exit(1)
	}

	return string(token)
}

func saveToken(token string) {
	err := os.WriteFile("token.txt", []byte(token), 0644)
	if err != nil {
		log.Fatalf("Error saving token")
		os.Exit(1)
	}
}
