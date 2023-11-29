package cmd

import (
	"fmt"
	"os"

	"github.com/doneill/er-cli-go/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// ----------------------------------------------
// user command
// ----------------------------------------------

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Current authenticated user data",
	Long:  `Return the currently authenticated er user data`,
	Run: func(cmd *cobra.Command, args []string) {
		user()
	},
}

// ----------------------------------------------
// functions
// ----------------------------------------------

func user() {
	userResponse, err := api.User()
	if err != nil {
		fmt.Println("Error authenticating:", err)
		os.Exit(1)
	}

	if userResponse != nil {
		data := []string{
			userResponse.Data.Username,
			userResponse.Data.Email,
			userResponse.Data.FirstName,
			userResponse.Data.LastName,
			userResponse.Data.ID,
			userResponse.Data.Pin,
			userResponse.Data.Subject.ID,
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Username", "Email", "First Name", "Last Name", "ID", "Pin", "Subject ID"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.Append(data)

		table.Render()
	}
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(userCmd)
}
