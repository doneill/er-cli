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
		name            string
		days            int
		status          string
		mockResponse    string
		expectedError   bool
		validateResult  func(*testing.T, *PatrolsResponse)
		validateRequest func(*testing.T, *http.Request)
	}{
		{
			name:   "successful response without filters",
			days:   0,
			status: "",
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
                                    "id": "segment123",
                                    "leader": {"name": "John Doe"},
                                    "patrol_type": "boat_patrol",
                                    "start_location": {"latitude": 1.234, "longitude": 5.678},
                                    "end_location": null,
                                    "time_range": {
                                        "start_time": "2025-01-15T10:00:00.000Z",
                                        "end_time": null
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
			validateRequest: func(t *testing.T, r *http.Request) {
				if !strings.Contains(r.URL.String(), "exclude_empty_patrols=true") {
					t.Error("Expected exclude_empty_patrols parameter")
				}
				if strings.Contains(r.URL.String(), "status=") {
					t.Error("Unexpected status parameter")
				}
				if strings.Contains(r.URL.String(), "filter=") {
					t.Error("Unexpected filter parameter")
				}
			},
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
				segment := patrol.PatrolSegments[0]
				if segment.ID != "segment123" {
					t.Errorf("Expected segment ID 'segment123', got '%s'", segment.ID)
				}
				if segment.Leader == nil {
					t.Fatal("Expected non-nil leader")
				}
				if segment.Leader.Name != "John Doe" {
					t.Errorf("Expected leader name 'John Doe', got '%s'", segment.Leader.Name)
				}
				// Open patrol should not have end location or end time
				if segment.EndLocation != nil {
					t.Error("Expected nil end location for open patrol")
				}
				if segment.TimeRange.EndTime != nil {
					t.Error("Expected nil end time for open patrol")
				}
			},
		},
		{
			name:   "successful response with done patrol having end location and time",
			days:   7,
			status: "",
			mockResponse: `{
                "data": {
                    "count": 1,
                    "results": [
                        {
                            "id": "test456",
                            "serial_number": 1002,
                            "state": "done",
                            "title": "Completed Patrol",
                            "patrol_segments": [
                                {
                                    "id": "segment456",
                                    "leader": {"name": "Jane Smith"},
                                    "patrol_type": "foot_patrol",
                                    "start_location": {"latitude": 2.345, "longitude": 6.789},
                                    "end_location": {"latitude": 2.355, "longitude": 6.799},
                                    "time_range": {
                                        "start_time": "2025-01-15T10:00:00.000Z",
                                        "end_time": "2025-01-15T14:00:00.000Z"
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
			validateRequest: func(t *testing.T, r *http.Request) {
				if !strings.Contains(r.URL.String(), "exclude_empty_patrols=true") {
					t.Error("Expected exclude_empty_patrols parameter")
				}
				if !strings.Contains(r.URL.String(), "filter=") {
					t.Error("Expected filter parameter")
				}
				if !strings.Contains(r.URL.String(), "patrols_overlap_daterange") {
					t.Error("Expected patrols_overlap_daterange in filter")
				}
			},
			validateResult: func(t *testing.T, response *PatrolsResponse) {
				if response == nil {
					t.Fatal("Expected non-nil response")
				}
				if len(response.Data.Results) != 1 {
					t.Errorf("Expected 1 result, got %d", len(response.Data.Results))
				}
				patrol := response.Data.Results[0]
				if patrol.State != "done" {
					t.Errorf("Expected state 'done', got '%s'", patrol.State)
				}
				if len(patrol.PatrolSegments) == 0 {
					t.Fatal("Expected at least one patrol segment")
				}
				segment := patrol.PatrolSegments[0]
				if segment.ID != "segment456" {
					t.Errorf("Expected segment ID 'segment456', got '%s'", segment.ID)
				}
				// Done patrol should have end location and end time
				if segment.EndLocation == nil {
					t.Fatal("Expected non-nil end location for done patrol")
				}
				if segment.EndLocation.Latitude != 2.355 {
					t.Errorf("Expected end latitude 2.355, got %f", segment.EndLocation.Latitude)
				}
				if segment.EndLocation.Longitude != 6.799 {
					t.Errorf("Expected end longitude 6.799, got %f", segment.EndLocation.Longitude)
				}
				if segment.TimeRange.EndTime == nil {
					t.Fatal("Expected non-nil end time for done patrol")
				}
				if *segment.TimeRange.EndTime != "2025-01-15T14:00:00.000Z" {
					t.Errorf("Expected end time '2025-01-15T14:00:00.000Z', got '%s'", *segment.TimeRange.EndTime)
				}
			},
		},
		{
			name:   "successful response with status filter",
			days:   7,
			status: "active",
			mockResponse: `{
                "data": {
                    "count": 1,
                    "results": [
                        {
                            "id": "test789",
                            "serial_number": 1003,
                            "state": "active",
                            "patrol_segments": [
                                {
                                    "id": "segment789",
                                    "patrol_type": "vehicle_patrol",
                                    "start_location": {"latitude": 3.456, "longitude": 7.890},
                                    "end_location": null,
                                    "time_range": {
                                        "start_time": "2025-01-15T08:00:00.000Z",
                                        "end_time": null
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
			validateRequest: func(t *testing.T, r *http.Request) {
				if !strings.Contains(r.URL.String(), "exclude_empty_patrols=true") {
					t.Error("Expected exclude_empty_patrols parameter")
				}
				if !strings.Contains(r.URL.String(), "status=active") {
					t.Error("Expected status parameter")
				}
				if !strings.Contains(r.URL.String(), "filter=") {
					t.Error("Expected filter parameter")
				}
			},
			validateResult: func(t *testing.T, response *PatrolsResponse) {
				if response == nil {
					t.Fatal("Expected non-nil response")
				}
				if len(response.Data.Results) != 1 {
					t.Errorf("Expected 1 result, got %d", len(response.Data.Results))
				}
				patrol := response.Data.Results[0]
				if patrol.State != "active" {
					t.Errorf("Expected state 'active', got '%s'", patrol.State)
				}
				if len(patrol.PatrolSegments) == 0 {
					t.Fatal("Expected at least one patrol segment")
				}
				segment := patrol.PatrolSegments[0]
				if segment.ID != "segment789" {
					t.Errorf("Expected segment ID 'segment789', got '%s'", segment.ID)
				}
				// Active patrol should not have end location or end time
				if segment.EndLocation != nil {
					t.Error("Expected nil end location for active patrol")
				}
				if segment.TimeRange.EndTime != nil {
					t.Error("Expected nil end time for active patrol")
				}
			},
		},
		{
			name:          "error response",
			days:          0,
			status:        "",
			mockResponse:  `{"status": {"code": 500, "message": "Internal Server Error"}}`,
			expectedError: true,
			validateRequest: func(t *testing.T, r *http.Request) {
				if !strings.Contains(r.URL.String(), "exclude_empty_patrols=true") {
					t.Error("Expected exclude_empty_patrols parameter")
				}
			},
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

				if tt.validateRequest != nil {
					tt.validateRequest(t, r)
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
			response, err := client.Patrols(tt.days, tt.status)

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
	_, err := client.Patrols(7, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
