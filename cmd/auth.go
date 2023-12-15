package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/doneill/er-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
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
// functions
// ----------------------------------------------

func auth() {
	fmt.Println("Enter password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	passwordStr := strings.TrimSpace(string(password))

	authResponse, err := api.Authenticate(SITENAME, USERNAME, passwordStr)
	if err != nil {
		fmt.Println("Error authenticating:", err)
		os.Exit(1)
	}

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
