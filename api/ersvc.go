package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ----------------------------------------------
// const endpoints
// ----------------------------------------------
const DOMAIN = ".pamdas.org"

const API_V1 = "/api/v1.0"

const API_AUTH = "/oauth2/token"

const API_USER = API_V1 + "/user"

const API_USER_ME = API_USER + "/me"

const API_USER_PROFILES = "/profiles"

// ----------------------------------------------
// Struct
// ----------------------------------------------

type Client struct {
	httpClient *http.Client
	sitename   string
	token      string
}

// ----------------------------------------------
// Functions
// ----------------------------------------------

func ERClient(sitename, token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		sitename: sitename,
		token:    token,
	}
}

func (c *Client) newRequest(method, endpoint string, isAuth bool) (*http.Request, error) {
	url := getApiUrl(c.sitename, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if isAuth {
		// Auth request headers
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")
	} else {
		// Regular API request headers
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Cache-control", "no-cache")
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// ----------------------------------------------
// Helper functions
// ----------------------------------------------

func getApiUrl(sitename string, endpoint string) string {
	return fmt.Sprintf("https://%s%s%s", sitename, DOMAIN, endpoint)
}

func getAuthRequest(sitename string) (*http.Request, error) {
	client := ERClient(sitename, "")
	return client.newRequest("POST", API_AUTH, true)
}

func getClientRequest(sitename string, endpoint string, token string) (*http.Request, error) {
	client := ERClient(sitename, token)
	return client.newRequest("GET", endpoint, false)
}
