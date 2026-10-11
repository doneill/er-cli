package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoRequestUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := ERClient("test", "expired-token", server.URL)
	_, err := client.User()

	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestDoRequestOtherStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := ERClient("test", "test-token", server.URL)
	_, err := client.User()

	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected a non-401 status not to be ErrUnauthorized, got %v", err)
	}
}
