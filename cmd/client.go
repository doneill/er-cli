package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/doneill/er-cli/api"
	"github.com/doneill/er-cli/config"
)

// ----------------------------------------------
// const var
// ----------------------------------------------

const tokenRefreshWindow = 24 * time.Hour

// ----------------------------------------------
// functions
// ----------------------------------------------

func newClient() *api.Client {
	sitename := config.Sitename()

	if config.NeedsRefresh(config.ExpiresAt(), time.Now(), tokenRefreshWindow) {
		if err := refreshToken(sitename); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not refresh token: %v\n", err)
		}
	}

	return api.ERClient(sitename, config.Token())
}

func refreshToken(sitename string) error {
	refresh := config.RefreshToken()
	if refresh == "" {
		return fmt.Errorf("no refresh token saved, run 'er auth' to sign in again")
	}

	authResponse, err := api.RefreshToken(sitename, refresh)
	if err != nil {
		return err
	}

	return saveTokens(authResponse)
}
