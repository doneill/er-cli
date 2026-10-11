package config

import (
	"testing"
	"time"
)

func TestNeedsRefresh(t *testing.T) {
	now := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	window := 24 * time.Hour

	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{"unknown expiry", time.Time{}, false},
		{"expires in 10 days", now.Add(10 * 24 * time.Hour), false},
		{"expires just outside window", now.Add(window + time.Minute), false},
		{"expires just inside window", now.Add(window - time.Minute), true},
		{"expires in 2 hours", now.Add(2 * time.Hour), true},
		{"already expired", now.Add(-time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NeedsRefresh(tt.expiresAt, now, window); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
