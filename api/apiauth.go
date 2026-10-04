package api

import (
	"fmt"
	"net/url"
	"strings"
)

// ----------------------------------------------
// structs
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
// exported functions
// ----------------------------------------------
func Authenticate(sitename, username, password string, opts ...string) (*AuthResponse, error) {
	client := ERClient(sitename, "", opts...)

	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	form.Set("client_id", "er_mobile_tracker")
	form.Set("grant_type", "password")

	req, err := client.newRequest("POST", API_AUTH, strings.NewReader(form.Encode()), true)
	if err != nil {
		return nil, fmt.Errorf("error generating auth request: %w", err)
	}

	var responseData AuthResponse
	if err := client.doRequest(req, &responseData); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &responseData, nil
}
