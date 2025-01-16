package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
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
		Name string `json:"name"`
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

type DateRangeFilter struct {
	DateRange struct {
		Lower string `json:"lower"`
		Upper string `json:"upper"`
	} `json:"date_range"`
	PatrolsOverlapDaterange bool `json:"patrols_overlap_daterange"`
}

// ----------------------------------------------
// Client methods
// ----------------------------------------------

func (c *Client) Patrols(days int) (*PatrolsResponse, error) {
	var endpoint string

	if days > 0 {
		now := time.Now().UTC()
		upper := now
		lower := now.AddDate(0, 0, -days)

		filter := DateRangeFilter{
			PatrolsOverlapDaterange: false,
		}
		filter.DateRange.Lower = lower.Format("2006-01-02T15:04:05.000Z")
		filter.DateRange.Upper = upper.Format("2006-01-02T15:04:05.000Z")

		filterJSON, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal date filter: %w", err)
		}

		params := url.Values{}
		params.Add("filter", string(filterJSON))
		params.Add("exclude_empty_patrols", "true")

		endpoint = fmt.Sprintf("%s?%s", API_PATROLS, params.Encode())
	} else {
		endpoint = fmt.Sprintf("%s?exclude_empty_patrols=true", API_PATROLS)
	}

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
