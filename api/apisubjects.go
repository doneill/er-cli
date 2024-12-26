package api

import (
	"fmt"
	"net/url"
)

// ----------------------------------------------
// structs
// ----------------------------------------------

type SubjectsResponse struct {
	Data   []Subject `json:"data"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

type Subject struct {
	ID               string       `json:"id"`
	LastPosition     LastPosition `json:"last_position"`
	LastPositionDate string       `json:"last_position_date"`
	Name             string       `json:"name"`
	SubjectSubtype   string       `json:"subject_subtype"`
	SubjectType      string       `json:"subject_type"`
}

type LastPosition struct {
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
	Type       string     `json:"type"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type Properties struct {
	DateTime             string     `json:"DateTime"`
	CoordinateProperties Coordinate `json:"coordinateProperties"`
	ID                   string     `json:"id"`
	Image                string     `json:"image"`
	LastVoiceCallStartAt *string    `json:"last_voice_call_start_at"`
	LocationRequestedAt  *string    `json:"location_requested_at"`
	RadioState           string     `json:"radio_state"`
	RadioStateAt         string     `json:"radio_state_at"`
	Stroke               string     `json:"stroke"`
	StrokeOpacity        float64    `json:"stroke-opacity"`
	StrokeWidth          int        `json:"stroke-width"`
	SubjectSubtype       string     `json:"subject_subtype"`
	SubjectType          string     `json:"subject_type"`
	Title                string     `json:"title"`
}

type Coordinate struct {
	Time string `json:"time"`
}

// ----------------------------------------------
// Client methods
// ----------------------------------------------

func (c *Client) Subjects(updatedSince string) (*SubjectsResponse, error) {
	params := url.Values{}
	params.Add("updated_since", updatedSince)

	endpoint := fmt.Sprintf("%s?%s", API_SUBJECTS, params.Encode())

	req, err := c.newRequest("GET", endpoint, false)
	if err != nil {
		return nil, fmt.Errorf("error generating subjects request: %w", err)
	}

	var responseData SubjectsResponse
	if err := c.doRequest(req, &responseData); err != nil {
		return nil, fmt.Errorf("error fetching subjects: %w", err)
	}

	return &responseData, nil
}
