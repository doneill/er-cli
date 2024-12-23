package api

import (
	"fmt"

	"github.com/doneill/er-cli/config"
)

// ----------------------------------------------
// stucts
// ----------------------------------------------

type UserResponse struct {
	Data struct {
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
	} `json:"data"`
	ErrorDescription string `json:"error_description"`
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
