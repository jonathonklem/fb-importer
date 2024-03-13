package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Printf("Choose either fetch or refresh, e.g.:.\n%s fetch\n-or-\n%s refresh\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "fetch" {
		fetch()
	} else {
		refreshToken()
	}
}
