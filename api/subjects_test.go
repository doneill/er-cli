package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSubjects(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse *SubjectsResponse
		statusCode     int
		expectError    bool
	}{
		{
			name: "successful response",
			serverResponse: &SubjectsResponse{
				Data: struct {
					Count    int       `json:"count"`
					Next     *string   `json:"next"`
					Previous *string   `json:"previous"`
					Results  []Subject `json:"results"`
				}{
					Count: 1,
					Results: []Subject{
						{
							ID: "123778d5-ffcc-4911-8d3b-e43cfdb426f7",
							LastPosition: LastPosition{
								Geometry: Geometry{
									Coordinates: []float64{-121.6670876888658, 47.44309785582009},
									Type:        "Point",
								},
								Type: "Feature",
							},
							LastPositionDate: "2024-12-23T18:34:51+00:00",
							Name:             "Test Subject",
							SubjectSubtype:   "ranger",
							SubjectType:      "person",
						},
					},
				},
				Status: struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				}{
					Code:    200,
					Message: "OK",
				},
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:           "server error",
			serverResponse: nil,
			statusCode:     http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if r.Header.Get("Authorization") != "Bearer testtoken" {
					t.Errorf("expected Bearer testtoken, got %s", r.Header.Get("Authorization"))
				}

				query := r.URL.Query()
				updatedSince := query.Get("updated_since")
				if updatedSince == "" {
					t.Error("expected updated_since parameter, got none")
				}

				updatedSinceTime, err := time.Parse("2006-01-02T15:04:05.000", updatedSince)
				if err != nil {
					t.Errorf("invalid date format: %v", err)
				}
				expectedTime := time.Now().AddDate(0, 0, -3)
				if updatedSinceTime.Day() != expectedTime.Day() {
					t.Errorf("expected date around %v, got %v", expectedTime, updatedSinceTime)
				}

				w.WriteHeader(tt.statusCode)
				if tt.serverResponse != nil {
					if err := json.NewEncoder(w).Encode(tt.serverResponse); err != nil {
						t.Errorf("failed to encode response: %v", err)
						return
					}
				}
			}))
			defer server.Close()

			client := ERClient("test", "testtoken", server.URL)
			resp, err := client.Subjects()

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if len(resp.Data.Results) != len(tt.serverResponse.Data.Results) {
					t.Errorf("expected %d subjects, got %d",
						len(tt.serverResponse.Data.Results), len(resp.Data.Results))
				}
				if len(resp.Data.Results) > 0 {
					if resp.Data.Results[0].ID != tt.serverResponse.Data.Results[0].ID {
						t.Errorf("expected subject ID %s, got %s",
							tt.serverResponse.Data.Results[0].ID, resp.Data.Results[0].ID)
					}
				}
			}
		})
	}
}
