package cmd

import (
	"strings"
	"testing"

	"github.com/doneill/er-cli/api"
)

func TestMaskPin(t *testing.T) {
	tests := []struct {
		name     string
		pin      string
		show     bool
		expected string
	}{
		{"masked by default", "1234", false, "****"},
		{"mask hides length", "12345678", false, "****"},
		{"shown with flag", "1234", true, "1234"},
		{"empty stays empty when masked", "", false, ""},
		{"empty stays empty when shown", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskPin(tt.pin, tt.show); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestUserOutputMasksPin(t *testing.T) {
	user := api.UserData{Username: "ranger", Pin: "1234"}

	tests := []struct {
		name     string
		show     bool
		expected string
	}{
		{"masked by default", false, pinMask},
		{"shown with --show-pin", true, "1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := showPin
			showPin = tt.show
			t.Cleanup(func() { showPin = original })

			// Pin is the sixth column: Username, Email, First Name, Last Name, ID, Pin, Subject ID
			row := formatTableData(&api.UserResponse{Data: user})
			if row[5] != tt.expected {
				t.Errorf("table: expected pin column %q, got %q", tt.expected, row[5])
			}

			if got := formatAllUserData(&api.UserResponse{Data: user}); !strings.Contains(got, "pin: "+tt.expected+"\n") {
				t.Errorf("all user data: expected %q in output, got:\n%s", "pin: "+tt.expected, got)
			}

			if got := formatAllProfileData(&user); !strings.Contains(got, "pin: "+tt.expected+"\n") {
				t.Errorf("all profile data: expected %q in output, got:\n%s", "pin: "+tt.expected, got)
			}
		})
	}
}
