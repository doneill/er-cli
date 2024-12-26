package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"time"
)

func TestSubjects(t *testing.T) {
	tests := []struct {
		name           string
		updatedSince   string
		serverResponse *SubjectsResponse
		statusCode     int
		expectError    bool
	}{
		{
			name:         "successful response",
			updatedSince: time.Now().AddDate(0, 0, -3).UTC().Format("2006-01-02T15:04:05.000"),
			serverResponse: &SubjectsResponse{
				Data: []Subject{
					{
						ID:   "123778d5-ffcc-4911-8d3b-e43cfdb426f7",
						Name: "Test Subject",
						LastPosition: LastPosition{
							Geometry: Geometry{
								Coordinates: []float64{-121.6670876888658, 47.44309785582009},
								Type:        "Point",
							},
							Properties: Properties{
								DateTime: "2024-12-23T18:34:51+00:00",
								Title:    "Test Subject",
							},
							Type: "Feature",
						},
						LastPositionDate: "2024-12-23T18:34:51+00:00",
						SubjectType:      "person",
						SubjectSubtype:   "ranger",
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
			updatedSince:   time.Now().AddDate(0, 0, -3).UTC().Format("2006-01-02T15:04:05.000"),
			serverResponse: nil,
			statusCode:     http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:         "empty response",
			updatedSince: time.Now().AddDate(0, 0, -3).UTC().Format("2006-01-02T15:04:05.000"),
			serverResponse: &SubjectsResponse{
				Data: []Subject{},
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

				if !strings.HasSuffix(r.URL.Path, API_SUBJECTS) {
					t.Errorf("expected path to end with %s, got %s", API_SUBJECTS, r.URL.Path)
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
			resp, err := client.Subjects(tt.updatedSince)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && tt.serverResponse != nil {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}

				if len(resp.Data) != len(tt.serverResponse.Data) {
					t.Errorf("expected %d subjects, got %d",
						len(tt.serverResponse.Data), len(resp.Data))
				}

				if len(resp.Data) > 0 {
					expectedSubject := tt.serverResponse.Data[0]
					actualSubject := resp.Data[0]

					if actualSubject.ID != expectedSubject.ID {
						t.Errorf("expected subject ID %s, got %s",
							expectedSubject.ID, actualSubject.ID)
					}

					if actualSubject.Name != expectedSubject.Name {
						t.Errorf("expected subject name %s, got %s",
							expectedSubject.Name, actualSubject.Name)
					}

					if len(actualSubject.LastPosition.Geometry.Coordinates) != 2 {
						t.Error("expected coordinates to have latitude and longitude")
					}
				}

				if resp.Status.Code != tt.serverResponse.Status.Code {
					t.Errorf("expected status code %d, got %d",
						tt.serverResponse.Status.Code, resp.Status.Code)
				}
			}
		})
	}
}

func TestSubjectTracks(t *testing.T) {
	tests := []struct {
		name           string
		subjectID      string
		daysAgo        int
		serverResponse *TracksResponse
		statusCode     int
		expectError    bool
	}{
		{
			name:      "successful response",
			subjectID: "123778d5-ffcc-4911-8d3b-e43cfdb426f7",
			daysAgo:   3,
			serverResponse: &TracksResponse{
				Data: struct {
					Features []Feature `json:"features"`
					Type     string    `json:"type"`
				}{
					Features: []Feature{
						{
							Geometry: TrackGeometry{
								Coordinates: [][]float64{
									{-121.6670876888658, 47.44309785582009},
									{-121.6695629568987, 47.44295692585065},
								},
								Type: "LineString",
							},
							Properties: TrackProperties{
								CoordinateProperties: struct {
									Times []string `json:"times"`
								}{
									Times: []string{
										"2024-12-23T18:34:51+00:00",
										"2024-12-23T18:34:45+00:00",
									},
								},
								ID:             "123778d5-ffcc-4911-8d3b-e43cfdb426f7",
								Title:          "Test Subject",
								SubjectType:    "person",
								SubjectSubtype: "ranger",
							},
							Type: "Feature",
						},
					},
					Type: "FeatureCollection",
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
			name:           "subject not found",
			subjectID:      "nonexistent",
			daysAgo:        3,
			serverResponse: nil,
			statusCode:     http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "empty subject ID",
			subjectID:      "",
			daysAgo:        3,
			serverResponse: nil,
			statusCode:     http.StatusBadRequest,
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

				expectedPath := path.Join(API_SUBJECT, tt.subjectID, API_SUBJECT_TRACKS)
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
				}

				query := r.URL.Query()
				since := query.Get("since")
				if since == "" {
					t.Error("expected since parameter, got none")
				}

				_, err := time.Parse("2006-01-02T15:04:05+00:00", since)
				if err != nil {
					t.Errorf("invalid date format: %v", err)
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
			resp, err := client.SubjectTracks(tt.subjectID, tt.daysAgo)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && tt.serverResponse != nil {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}

				if len(resp.Data.Features) != len(tt.serverResponse.Data.Features) {
					t.Errorf("expected %d features, got %d",
						len(tt.serverResponse.Data.Features), len(resp.Data.Features))
				}

				if len(resp.Data.Features) > 0 {
					if resp.Data.Features[0].Type != tt.serverResponse.Data.Features[0].Type {
						t.Errorf("expected feature type %s, got %s",
							tt.serverResponse.Data.Features[0].Type, resp.Data.Features[0].Type)
					}

					if len(resp.Data.Features[0].Geometry.Coordinates) !=
						len(tt.serverResponse.Data.Features[0].Geometry.Coordinates) {
						t.Error("coordinates length mismatch")
					}
				}
			}
		})
	}
}
