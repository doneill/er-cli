package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ----------------------------------------------
// const var
// ----------------------------------------------

const DOMAIN = ".pamdas.org"
const API_AUTH = "/oauth2/token"

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
// exported funtions
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

func Authenticate(sitename, username, password string) (*Response, error) {
	client := &http.Client{}
	authReq := getAuthRequest(sitename)
	authReq.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("username=%s&password=%s&client_id=er_mobile_tracker&grant_type=password", username, password)))

	res, err := client.Do(&authReq)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer res.Body.Close()

	var responseData Response
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
	}

	if res.StatusCode == 200 {
		return &responseData, nil
	}

	fmt.Println("Error:", res.StatusCode)
	fmt.Println("Error Description:", responseData.ErrorDescription)

	return nil, err
}
