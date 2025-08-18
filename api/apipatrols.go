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
	ID     string `json:"id"`
	Leader *struct {
		Name string `json:"name"`
	} `json:"leader"`
	PatrolType    string    `json:"patrol_type"`
	StartLocation *Location `json:"start_location"`
	EndLocation   *Location `json:"end_location"`
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

func (c *Client) Patrols(days int, status string) (*PatrolsResponse, error) {
	params := url.Values{}
	params.Add("exclude_empty_patrols", "true")
	params.Add("page_size", "200")

	if status != "" {
		params.Add("status", status)
	}

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

		params.Add("filter", string(filterJSON))
	}

	endpoint := fmt.Sprintf("%s?%s", API_PATROLS, params.Encode())

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
func (c *Client) PatrolTracks(subjectID string, since string, until string) (*TracksResponse, error) {
	params := url.Values{}
	params.Add("since", since)
	params.Add("until", until)
	endpoint := path.Join(API_SUBJECT, subjectID, API_SUBJECT_TRACKS)
	endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	req, err := c.newRequest("GET", endpoint, false)
	if err != nil {
		return nil, fmt.Errorf("error generating patrol tracks request: %w", err)
	}
	var responseData TracksResponse
	if err := c.doRequest(req, &responseData); err != nil {
		return nil, fmt.Errorf("error fetching patrol tracks: %w", err)
	}
}

	return &responseData, nil
