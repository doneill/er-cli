package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// ----------------------------------------------
// exported funtions
// ----------------------------------------------

func Sitename() string {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Error: er not configured properly, try reauthenticating")
		} else {
			fmt.Println("Error:", err)
		}
	} else {
		return viper.Get("sitename").(string)
	}
	return ""
}

func ExpiresAt() time.Time {
	if err := viper.ReadInConfig(); err != nil {
		return time.Time{}
	}

	seconds := viper.GetInt64("expires_at")
	if seconds == 0 {
		return time.Time{}
	}

	return time.Unix(seconds, 0)
}

func NeedsRefresh(expiresAt time.Time, now time.Time, window time.Duration) bool {
	if expiresAt.IsZero() {
		return false
	}
	return now.Add(window).After(expiresAt)
}

func RefreshToken() string {
	if err := viper.ReadInConfig(); err != nil {
		return ""
	}
	return viper.GetString("refresh_token")
}

func Token() string {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Error: er not configured properly, try reauthenticating")
		} else {
			fmt.Println("Error:", err)
		}
	} else {
		return viper.Get("oauth_token").(string)
	}
	return ""
}
