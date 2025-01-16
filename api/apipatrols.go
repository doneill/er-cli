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
	PatrolType    string    `json:"patrol_type"`
	StartLocation *Location `json:"start_location"`
	TimeRange     TimeRange `json:"time_range"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TimeRange struct {
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
}

// ----------------------------------------------
// Client methods
// ----------------------------------------------

func (c *Client) Patrols() (*PatrolsResponse, error) {
	endpoint := fmt.Sprintf("%s?exclude_empty_patrols=true", API_PATROLS)

	req, err := c.newRequest("GET", endpoint, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create patrols request: %w", err)
	}

	var response PatrolsResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, fmt.Errorf("failed to get patrols: %w", err)
	}

	return &response, nil
}
