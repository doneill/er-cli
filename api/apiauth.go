package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ----------------------------------------------
// stucts
// ----------------------------------------------

type AuthResponse struct {
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

func Authenticate(sitename, username, password string) (*AuthResponse, error) {
	client := &http.Client{}
	authReq, err := getAuthRequest(sitename)
	if err != nil {
		fmt.Println("Error generating auth client", err)
	}
	authReq.Body = io.NopCloser(
		strings.NewReader(
			fmt.Sprintf(
				"username=%s&password=%s&client_id=er_mobile_tracker&grant_type=password", username, password)))

	res, err := client.Do(authReq)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer res.Body.Close()

	var responseData AuthResponse
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
