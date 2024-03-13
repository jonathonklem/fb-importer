package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func refreshToken() {
	current_token := getToken()
	url := os.Getenv("FB_REFRESH_URL") + "&fb_exchange_token=" + current_token

	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading body.\n", err)
	}

	var data TokenResponse

	err = json.Unmarshal([]byte(body), &data)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Setenv("ACCESS_TOKEN", data.AccessToken)
	saveToken(data.AccessToken)
}
