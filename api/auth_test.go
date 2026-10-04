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
