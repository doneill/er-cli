package api

import (
	"fmt"
	"net/http"
)

// ----------------------------------------------
// const endpoints
// ----------------------------------------------

const DOMAIN = ".pamdas.org"

const API_V1 = "/api/v1.0"

const API_AUTH = "/oauth2/token"

const API_USER = API_V1 + "/user"

const API_USER_ME = API_USER + "/me"

// ----------------------------------------------
// funtions
// ----------------------------------------------

func getApiUrl(sitename string, endpoint string) string {
	return fmt.Sprintf("https://%s%s%s", sitename, DOMAIN, endpoint)
}

func getAuthRequest(sitename string) http.Request {
	req, err := http.NewRequest("POST", getApiUrl(sitename, API_AUTH), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return *req
}

func getClientRequest(sitename string, endpoint string, token string) http.Request {
	req, err := http.NewRequest("GET", getApiUrl(sitename, endpoint), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Cache-control", "no-cache")

	return *req
}
