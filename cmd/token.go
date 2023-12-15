package cmd

import (
	"fmt"

	"github.com/doneill/er-cli/config"
	"github.com/spf13/cobra"
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
	var token = config.Token()
	fmt.Println(token)
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	authCmd.AddCommand(tokenCmd)
}
