package cmd

import (
	"fmt"
	"log"
	"strings"
	"syscall"
	"time"

	"github.com/doneill/er-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
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

	viper.Set("user", USERNAME)
	viper.Set("sitename", SITENAME)

	if err := saveTokens(authResponse); err != nil {
		log.Fatalf("Error updating configuration: %v", err)
	}

	fmt.Println("Authenticated!")
}

func saveTokens(authResponse *api.AuthResponse) error {
	expiresAt := time.Now().Add(time.Duration(authResponse.ExpiresIn) * time.Second)

	viper.Set("oauth_token", authResponse.AccessToken)
	viper.Set("refresh_token", authResponse.RefreshToken)
	viper.Set("expires_at", expiresAt.Unix())

	viper.SetConfigPermissions(0o600)
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
