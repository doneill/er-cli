package api

import (
	"fmt"
	"path"

	"github.com/doneill/er-cli/config"
)

// ----------------------------------------------
// stucts
// ----------------------------------------------

type UserData struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Role        string `json:"role"`
	IsStaff     bool   `json:"is_staff"`
	IsSuperUser bool   `json:"is_superuser"`
	DateJoined  string `json:"date_joined"`
	ID          string `json:"id"`
	IsActive    bool   `json:"is_active"`
	LastLogin   string `json:"last_login"`
	Pin         string `json:"pin"`
	Subject     struct {
		ID string `json:"id"`
	} `json:"subject"`
	AcceptedEula bool `json:"accepted_eula"`
}

type UserProfilesResponse struct {
	Data   []UserData `json:"data"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

type UserResponse struct {
	Data             UserData `json:"data"`
	ErrorDescription string   `json:"error_description"`
}

// ----------------------------------------------
// exported funtions
// ----------------------------------------------

func User() (*UserResponse, error) {
	client := ERClient(config.Sitename(), config.Token())

	req, err := client.newRequest("GET", API_USER_ME, false)
	if err != nil {
		return nil, fmt.Errorf("error generating request: %w", err)
	}

	var responseData UserResponse
	if err := client.doRequest(req, &responseData); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return &responseData, nil
}

func UserProfiles(userID string) (*UserProfilesResponse, error) {
	client := ERClient(config.Sitename(), config.Token())

	// Construct the profiles endpoint
	profilesEndpoint := path.Join(API_USER, userID, API_USER_PROFILES)

	req, err := client.newRequest("GET", profilesEndpoint, false)
	if err != nil {
		return nil, fmt.Errorf("error generating profiles request: %w", err)
	}

	var responseData UserProfilesResponse
	if err := client.doRequest(req, &responseData); err != nil {
		return nil, fmt.Errorf("error fetching profiles: %w", err)
	}

	return &responseData, nil
}
