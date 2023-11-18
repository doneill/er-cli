package main

import (
	"fmt"
	"os"
)

type Response struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	RefreshToken     string `json:"refresh_token"`
	ErrorDescription string `json:"error_description"`
}

func main() {
	// Get the sitename, username, and password from CLI arguments
	fmt.Println("Enter sitename:")
	var sitename string
	_, err := fmt.Scan(&sitename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Enter username:")
	var username string
	_, err = fmt.Scan(&username)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Enter password:")
	var password string
	_, err = fmt.Scan(&password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Call the authenticate function to get the access token and expires in
	response, err := authenticate(sitename, username, password)
	if err != nil {
		fmt.Println("Error authenticating:", err)
		os.Exit(1)
	}

	// Print out the access token and expires in if the request was successful
	if response != nil {
		fmt.Printf("Access Token: %s\n", response.AccessToken)
		fmt.Printf("Expires In: %d\n", response.ExpiresIn)
	}
}
