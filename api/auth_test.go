package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		statusCode  int
		expectError bool
	}{
		{
			name:        "successful response",
			username:    "testuser",
			password:    "password",
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "special characters are form encoded",
			username:    "test+user@example.com",
			password:    "p&ss=1+2% é",
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "invalid credentials",
			username:    "testuser",
			password:    "wrong",
			statusCode:  http.StatusUnauthorized,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != API_AUTH {
					t.Errorf("expected path %s, got %s", API_AUTH, r.URL.Path)
				}
				if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
					t.Errorf("expected form content type, got %s", r.Header.Get("Content-Type"))
				}
				if r.ContentLength <= 0 {
					t.Errorf("expected Content-Length to be set, got %d", r.ContentLength)
				}
				if err := r.ParseForm(); err != nil {
					t.Fatalf("failed to parse form: %v", err)
				}

				expected := map[string]string{
					"username":   tt.username,
					"password":   tt.password,
					"client_id":  "er_mobile_tracker",
					"grant_type": "password",
				}
				for key, want := range expected {
					if got := r.PostForm.Get(key); got != want {
						t.Errorf("expected %s %q, got %q", key, want, got)
					}
				}
				if len(r.PostForm) != len(expected) {
					t.Errorf("expected %d form fields, got %d: %v", len(expected), len(r.PostForm), r.PostForm)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					if err := json.NewEncoder(w).Encode(AuthResponse{
						AccessToken: "testtoken",
						TokenType:   "Bearer",
					}); err != nil {
						t.Errorf("failed to encode response: %v", err)
					}
				}
			}))
			defer server.Close()

			response, err := Authenticate("test", tt.username, tt.password, server.URL)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if response.AccessToken != "testtoken" {
				t.Errorf("expected access token testtoken, got %s", response.AccessToken)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
		statusCode   int
		response     string
		expectError  bool
	}{
		{
			name:         "successful refresh",
			refreshToken: "old-refresh-token",
			statusCode:   http.StatusOK,
			response:     `{"access_token": "new-access-token", "refresh_token": "new-refresh-token", "expires_in": 2592000, "token_type": "Bearer"}`,
			expectError:  false,
		},
		{
			name:         "expired refresh token",
			refreshToken: "expired-refresh-token",
			statusCode:   http.StatusBadRequest,
			response:     `{"error": "invalid_grant"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != API_AUTH {
					t.Errorf("expected path %s, got %s", API_AUTH, r.URL.Path)
				}
				if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
					t.Errorf("expected form content type, got %s", r.Header.Get("Content-Type"))
				}
				if err := r.ParseForm(); err != nil {
					t.Fatalf("failed to parse form: %v", err)
				}

				expected := map[string]string{
					"refresh_token": tt.refreshToken,
					"grant_type":    "refresh_token",
					"client_id":     "er_mobile_tracker",
				}
				for key, want := range expected {
					if got := r.PostForm.Get(key); got != want {
						t.Errorf("expected %s %q, got %q", key, want, got)
					}
				}
				if len(r.PostForm) != len(expected) {
					t.Errorf("expected %d form fields, got %d: %v", len(expected), len(r.PostForm), r.PostForm)
				}

				w.WriteHeader(tt.statusCode)
				if _, err := w.Write([]byte(tt.response)); err != nil {
					t.Errorf("failed to write response: %v", err)
				}
			}))
			defer server.Close()

			response, err := RefreshToken("test", tt.refreshToken, server.URL)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if response.AccessToken != "new-access-token" {
				t.Errorf("expected access token new-access-token, got %s", response.AccessToken)
			}
			if response.RefreshToken != "new-refresh-token" {
				t.Errorf("expected rotated refresh token new-refresh-token, got %s", response.RefreshToken)
			}
			if response.ExpiresIn != 2592000 {
				t.Errorf("expected expires_in 2592000, got %d", response.ExpiresIn)
			}
		})
	}
}
