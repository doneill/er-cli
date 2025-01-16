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

const API_ACTIVITY = API_V1 + "/activity"

const API_PATROLS = API_ACTIVITY + "/patrols"

const API_SUBJECT = API_V1 + "/subject"

const API_SUBJECTS = API_V1 + "/subjects"

const API_SUBJECT_TRACKS = "/tracks"

const API_USER = API_V1 + "/user"

const API_USER_ME = API_USER + "/me"

const API_USER_PROFILES = "/profiles"

// ----------------------------------------------
// struct
// ----------------------------------------------

type Client struct {
	httpClient *http.Client
	sitename   string
	token      string
	mockURL    string
}

// ----------------------------------------------
// functions
// ----------------------------------------------

func ERClient(sitename, token string, opts ...string) *Client {
	mockURL := ""
	if len(opts) > 0 {
		mockURL = opts[0]
	}
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		sitename: sitename,
		token:    token,
		mockURL:  mockURL,
	}
}

func (c *Client) newRequest(method, endpoint string, isAuth bool) (*http.Request, error) {
	url := getApiUrl(c.sitename, endpoint, c.mockURL)
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

func getApiUrl(sitename string, endpoint string, mockURL string) string {
	if mockURL != "" {
		return mockURL + endpoint
	}
	return fmt.Sprintf("https://%s%s%s", sitename, DOMAIN, endpoint)
}
