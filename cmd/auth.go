package cmd

import (
	"fmt"
	"github.com/doneill/er-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"log"
	"strings"
	"syscall"
)

// ----------------------------------------------
// static var
// ----------------------------------------------
var (
	SITENAME string
	USERNAME string
)

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
	if err != nil {
		log.Fatalf("Error reading password: %v", err)
	}

	passwordStr := strings.TrimSpace(string(password))
	authResponse, err := api.Authenticate(SITENAME, USERNAME, passwordStr)
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}

	if err := updateAuthConfig(authResponse); err != nil {
		log.Fatalf("Error updating configuration: %v", err)
	}

	fmt.Println("Authenticated!")
}

func updateAuthConfig(authResponse *api.AuthResponse) error {
	viper.Set("user", USERNAME)
	viper.Set("sitename", SITENAME)
	viper.Set("oauth_token", authResponse.AccessToken)
	viper.Set("expires", authResponse.ExpiresIn)

	if err := viper.WriteConfigAs(PROGRAM_NAME + CONFIG_TYPE); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
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
