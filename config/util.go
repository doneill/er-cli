package config

import (
	"fmt"

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
