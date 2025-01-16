package api

import (
	"fmt"
)

// ----------------------------------------------
// Patrol types
// ----------------------------------------------

type PatrolsResponse struct {
	Data struct {
		Count    int      `json:"count"`
		Next     string   `json:"next"`
		Previous string   `json:"previous"`
		Results  []Patrol `json:"results"`
	} `json:"data"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

type Patrol struct {
	ID             string          `json:"id"`
	SerialNumber   int             `json:"serial_number"`
	State          string          `json:"state"`
	Title          *string         `json:"title"`
	PatrolSegments []PatrolSegment `json:"patrol_segments"`
}

type PatrolSegment struct {
	Leader *struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"leader"`
	PatrolType string `json:"patrol_type"`
}

// ----------------------------------------------
// Client methods
// ----------------------------------------------

func (c *Client) Patrols() (*PatrolsResponse, error) {
	req, err := c.newRequest("GET", API_PATROLS, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create patrols request: %w", err)
	}

	var response PatrolsResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, fmt.Errorf("failed to get patrols: %w", err)
	}

	return &response, nil
}
