package api

import (
	"fmt"
	"net/url"
	"path"
	"time"
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

type TracksResponse struct {
	Data struct {
		Features []Feature `json:"features"`
		Type     string    `json:"type"`
	} `json:"data"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

type Feature struct {
	Geometry   TrackGeometry   `json:"geometry"`
	Properties TrackProperties `json:"properties"`
	Type       string          `json:"type"`
}

type TrackGeometry struct {
	Coordinates [][]float64 `json:"coordinates"`
	Type        string      `json:"type"`
}

type TrackProperties struct {
	CoordinateProperties struct {
		Times []string `json:"times"`
	} `json:"coordinateProperties"`
	ID                   string  `json:"id"`
	Image                string  `json:"image"`
	LastVoiceCallStartAt *string `json:"last_voice_call_start_at"`
	LocationRequestedAt  *string `json:"location_requested_at"`
	RadioState           string  `json:"radio_state"`
	RadioStateAt         string  `json:"radio_state_at"`
	Stroke               string  `json:"stroke"`
	StrokeOpacity        float64 `json:"stroke-opacity"`
	StrokeWidth          int     `json:"stroke-width"`
	SubjectSubtype       string  `json:"subject_subtype"`
	SubjectType          string  `json:"subject_type"`
	Title                string  `json:"title"`
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

func (c *Client) SubjectTracks(subjectID string, daysAgo int) (*TracksResponse, error) {
	sinceDate := time.Now().AddDate(0, 0, -daysAgo).UTC().Format("2006-01-02T15:04:05+00:00")

	params := url.Values{}
	params.Add("since", sinceDate)

	endpoint := path.Join(API_SUBJECT, subjectID, API_SUBJECT_TRACKS)
	endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := c.newRequest("GET", endpoint, false)
	if err != nil {
		return nil, fmt.Errorf("error generating tracks request: %w", err)
	}

	var responseData TracksResponse
	if err := c.doRequest(req, &responseData); err != nil {
		return nil, fmt.Errorf("error fetching tracks: %w", err)
	}

	return &responseData, nil
}
