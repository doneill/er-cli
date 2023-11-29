package cmd

import (
	"fmt"
	"os"

	"github.com/doneill/er-cli-go/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ----------------------------------------------
// static var
// ----------------------------------------------

var SITENAME string
var USERNAME string

// ----------------------------------------------
// auth command
// ----------------------------------------------

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication with EarthRanger",
	Long:  `Authenticate er with EarthRanger`,
	Run: func(cmd *cobra.Command, args []string) {
		auth()
	},
}

// ----------------------------------------------
// funtions
// ----------------------------------------------

func auth() {
	fmt.Println("Enter password:")
	var password string
	_, err := fmt.Scan(&password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Call the authenticate function to get the access token and expires in
	authResponse, err := api.Authenticate(SITENAME, USERNAME, password)
	if err != nil {
		fmt.Println("Error authenticating:", err)
		os.Exit(1)
	}

	// Print out the access token and expires in if the request was successful
	if authResponse != nil {
		viper.Set("user", USERNAME)
		viper.Set("sitename", SITENAME)
		viper.Set("oauth_token", authResponse.AccessToken)
		viper.Set("expires", authResponse.ExpiresIn)
		err := viper.WriteConfigAs(PROGRAM_NAME + CONFIG_TYPE)
		if err != nil {
			fmt.Println("Error writing configuration file:", err)
		} else {
			fmt.Println("Authenticated!")
		}
	}
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVarP(&SITENAME, "sitename", "s", "", "EarthRanger site name")
	authCmd.Flags().StringVarP(&USERNAME, "username", "u", "", "EarthRanger user name")
	authCmd.MarkFlagRequired("sitename")
	authCmd.MarkFlagRequired("username")
}
