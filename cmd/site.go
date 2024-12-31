package cmd

import (
	"fmt"

	"github.com/doneill/er-cli/config"
	"github.com/spf13/cobra"
)

// ----------------------------------------------
// site command
// ----------------------------------------------

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Display site",
	Long:  `Return the site the current users is authenticated on`,
	Run: func(cmd *cobra.Command, args []string) {
		site()
	},
}

// ----------------------------------------------
// funtions
// ----------------------------------------------

func site() {
	var site = config.Sitename()
	fmt.Println(site)
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	userCmd.AddCommand(siteCmd)
}
