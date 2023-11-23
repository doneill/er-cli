package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ----------------------------------------------
// token command
// ----------------------------------------------

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display token",
	Long:  `Return the current users token`,
	Run: func(cmd *cobra.Command, args []string) {
		token()
	},
}

// ----------------------------------------------
// funtions
// ----------------------------------------------

func token() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Error: er not configured properly, try reauthenticating")
		} else {
			fmt.Println("Error:", err)
		}
	} else {
		var token = viper.Get("oauth_token")
		fmt.Println(token)
	}
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	authCmd.AddCommand(tokenCmd)
}
