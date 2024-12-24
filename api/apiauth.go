package api

import (
	"fmt"
	"io"
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
func Authenticate(sitename, username, password string) (*AuthResponse, error) {
	client := ERClient(sitename, "")

	req, err := client.newRequest("POST", API_AUTH, true)
	if err != nil {
		return nil, fmt.Errorf("error generating auth request: %w", err)
	}

	authBody := fmt.Sprintf(
		"username=%s&password=%s&client_id=er_mobile_tracker&grant_type=password",
		username, password,
	)
	req.Body = io.NopCloser(strings.NewReader(authBody))

	var responseData AuthResponse
	if err := client.doRequest(req, &responseData); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &responseData, nil
}
