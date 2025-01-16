package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPatrols(t *testing.T) {
	tests := []struct {
		name           string
		days           int
		mockResponse   string
		expectedError  bool
		validateResult func(*testing.T, *PatrolsResponse)
	}{
		{
			name: "successful response without date filter",
			days: 0,
			mockResponse: `{
                "data": {
                    "count": 1,
                    "next": null,
                    "previous": null,
                    "results": [
                        {
                            "id": "test123",
                            "serial_number": 1001,
                            "state": "open",
                            "title": "Test Patrol",
                            "patrol_segments": [
                                {
                                    "leader": {"name": "John Doe"},
                                    "patrol_type": "boat_patrol",
                                    "start_location": {"latitude": 1.234, "longitude": 5.678},
                                    "time_range": {
                                        "start_time": "2025-01-15T10:00:00.000Z",
                                        "end_time": "2025-01-15T11:00:00.000Z"
                                    }
                                }
                            ]
                        }
                    ]
                },
                "status": {
                    "code": 200,
                    "message": "OK"
                }
            }`,
			expectedError: false,
			validateResult: func(t *testing.T, response *PatrolsResponse) {
				if response == nil {
					t.Fatal("Expected non-nil response")
				}
				if len(response.Data.Results) != 1 {
					t.Errorf("Expected 1 result, got %d", len(response.Data.Results))
				}
				patrol := response.Data.Results[0]
				if patrol.ID != "test123" {
					t.Errorf("Expected ID 'test123', got '%s'", patrol.ID)
				}
				if patrol.SerialNumber != 1001 {
					t.Errorf("Expected serial number 1001, got %d", patrol.SerialNumber)
				}
				if len(patrol.PatrolSegments) == 0 {
					t.Fatal("Expected at least one patrol segment")
				}
				if patrol.PatrolSegments[0].Leader == nil {
					t.Fatal("Expected non-nil leader")
				}
				if patrol.PatrolSegments[0].Leader.Name != "John Doe" {
					t.Errorf("Expected leader name 'John Doe', got '%s'", patrol.PatrolSegments[0].Leader.Name)
				}
			},
		},
		{
			name: "successful response with date filter",
			days: 7,
			mockResponse: `{
                "data": {
                    "count": 1,
                    "results": [
                        {
                            "id": "test456",
                            "serial_number": 1002,
                            "state": "closed"
                        }
                    ]
                },
                "status": {
                    "code": 200,
                    "message": "OK"
                }
            }`,
			expectedError: false,
			validateResult: func(t *testing.T, response *PatrolsResponse) {
				if response == nil {
					t.Fatal("Expected non-nil response")
				}
				if len(response.Data.Results) != 1 {
					t.Errorf("Expected 1 result, got %d", len(response.Data.Results))
				}
			},
		},
		{
			name:           "error response",
			days:           0,
			mockResponse:   `{"status": {"code": 500, "message": "Internal Server Error"}}`,
			expectedError:  true,
			validateResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Validate request
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				if tt.days > 0 {
					if !strings.Contains(r.URL.String(), "filter=") {
						t.Error("Expected filter parameter in URL for date-filtered request")
					}
					if !strings.Contains(r.URL.String(), "patrols_overlap_daterange") {
						t.Error("Expected patrols_overlap_daterange in filter")
					}
				}

				// Return mock response
				w.Header().Set("Content-Type", "application/json")
				if strings.Contains(tt.mockResponse, `"code": 500`) {
					w.WriteHeader(http.StatusInternalServerError)
				}
				if _, err := w.Write([]byte(tt.mockResponse)); err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			}))
			defer server.Close()

			client := ERClient("test", "test-token", server.URL)
			response, err := client.Patrols(tt.days)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validateResult != nil {
				tt.validateResult(t, response)
			}
		})
	}
}

func TestDateRangeFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filterStr := r.URL.Query().Get("filter")
		if filterStr == "" {
			t.Error("Expected filter parameter in URL")
			return
		}

		var filter DateRangeFilter
		err := json.Unmarshal([]byte(filterStr), &filter)
		if err != nil {
			t.Errorf("Failed to parse filter JSON: %v", err)
			return
		}

		// Validate date format
		_, err = time.Parse(time.RFC3339, filter.DateRange.Lower)
		if err != nil {
			t.Errorf("Invalid lower date format: %v", err)
		}

		_, err = time.Parse(time.RFC3339, filter.DateRange.Upper)
		if err != nil {
			t.Errorf("Invalid upper date format: %v", err)
		}

		if filter.PatrolsOverlapDaterange {
			t.Error("Expected PatrolsOverlapDaterange to be false")
		}

		// Return a valid response
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"data":{"count":0,"results":[]},"status":{"code":200,"message":"OK"}}`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := ERClient("test", "test-token", server.URL)
	_, err := client.Patrols(7)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
