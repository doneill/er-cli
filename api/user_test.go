package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse *UserResponse
		statusCode     int
		expectError    bool
	}{
		{
			name: "successful response",
			serverResponse: &UserResponse{
				Data: UserData{
					Username:  "testuser",
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
					ID:        "123",
					Subject: struct {
						ID string `json:"id"`
					}{ID: "456"},
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

				w.WriteHeader(tt.statusCode)
				if tt.serverResponse != nil {
					if err := json.NewEncoder(w).Encode(tt.serverResponse); err != nil {
						t.Errorf("failed to encode response: %v", err)
						return
					}
				}
			}))
			defer server.Close()

			er := ERClient("test", "testtoken", server.URL)
			resp, err := er.User()

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
				if resp.Data.Username != tt.serverResponse.Data.Username {
					t.Errorf("expected username %s, got %s",
						tt.serverResponse.Data.Username, resp.Data.Username)
				}
			}
		})
	}
}

func TestUserProfiles(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		serverResponse *UserProfilesResponse
		statusCode     int
		expectError    bool
	}{
		{
			name:   "successful response",
			userID: "123",
			serverResponse: &UserProfilesResponse{
				Data: []UserData{
					{
						Username:  "profile1",
						Email:     "profile1@example.com",
						FirstName: "Profile",
						LastName:  "One",
						ID:        "456",
						Subject: struct {
							ID string `json:"id"`
						}{ID: "789"},
					},
				},
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:           "server error",
			userID:         "123",
			serverResponse: nil,
			statusCode:     http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:           "invalid user ID",
			userID:         "",
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

				expectedPath := path.Join(API_USER, tt.userID, API_USER_PROFILES)
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
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

			er := ERClient("test", "testtoken", server.URL)
			resp, err := er.UserProfiles(tt.userID)

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
				if len(resp.Data) != len(tt.serverResponse.Data) {
					t.Errorf("expected %d profiles, got %d",
						len(tt.serverResponse.Data), len(resp.Data))
				}
				if len(resp.Data) > 0 {
					if resp.Data[0].Username != tt.serverResponse.Data[0].Username {
						t.Errorf("expected username %s, got %s",
							tt.serverResponse.Data[0].Username, resp.Data[0].Username)
					}
				}
			}
		})
	}
}
