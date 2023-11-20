package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ----------------------------------------------
// static var
// ----------------------------------------------

var DOMAIN = ".pamdas.org"
var API_AUTH = "/oauth2/token"

// ----------------------------------------------
// stucts
// ----------------------------------------------

type Response struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	RefreshToken     string `json:"refresh_token"`
	ErrorDescription string `json:"error_description"`
}

// ----------------------------------------------
// package funtions
// ----------------------------------------------

func authenticate(sitename, username, password string) (*Response, error) {
	// Create a new HTTP client and make a POST request to the authentication endpoint
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s%s", sitename, DOMAIN, API_AUTH), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set the Content-Type and Accept headers on the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Add the username, password, client_id, and grant_type to the request body
	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("username=%s&password=%s&client_id=er_mobile_tracker&grant_type=password", username, password)))

	// Make the POST request and get the response
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer res.Body.Close()

	// Decode the JSON response from the authentication endpoint
	var responseData Response
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
	}

	// If the request was successful, return the access token and expires in
	if res.StatusCode == 200 {
		return &responseData, nil
	}

	// Print out the error and error description if the request was rejected
	fmt.Println("Error:", res.StatusCode)
	fmt.Println("Error Description:", responseData.ErrorDescription)

	return nil, err
}
