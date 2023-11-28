package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ----------------------------------------------
// exported funtions
// ----------------------------------------------

func Token() string {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Error: er not configured properly, try reauthenticating")
		} else {
			fmt.Println("Error:", err)
		}
	} else {
		var token = viper.Get("oauth_token")
		return token.(string)
	}
	return ""
}
