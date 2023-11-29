package cmd

import (
	"fmt"
	"os"

	"github.com/doneill/er-cli-go/api"
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
		formattedResponse := fmt.Sprintf("username: %s\nemail: %s\nfirst name: %s\nlast name: %s\nrole: %s\nis staff: %t\nis superuser: %t\ndate joined: %s\nid: %s\nisactive: %t\nlast login: %s\npin: %s\nsubject id: %s\npermissions:\n  patrol: %v\nmobile tests: %v",
			userResponse.Data.Username,
			userResponse.Data.Email,
			userResponse.Data.FirstName,
			userResponse.Data.LastName,
			userResponse.Data.Role,
			userResponse.Data.IsStaff,
			userResponse.Data.IsSuperUser,
			userResponse.Data.DateJoined,
			userResponse.Data.ID,
			userResponse.Data.IsActive,
			userResponse.Data.LastLogin,
			userResponse.Data.Pin,
			userResponse.Data.Subject.ID,
			userResponse.Data.Permissions.Patrol,
			userResponse.Data.Permissions.MobileTests)
		fmt.Println(formattedResponse)
	}
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(userCmd)
}
