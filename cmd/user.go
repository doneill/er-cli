package cmd

import (
	"fmt"
	"os"

	"github.com/doneill/er-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var all bool

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
		switch {
		case all:
			allUserDataFmt := fmt.Sprintf("username: %s\nemail: %s\nfirst name: %s\nlast name: %s\nrole: %s\nis staff: %t\nis superuser: %t\ndate joined: %s\nid: %s\nisactive: %t\nlast login: %s\npin: %s\nsubject id: %s",
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
				userResponse.Data.Subject.ID)
			fmt.Println(allUserDataFmt)
		default:
			userData := []string{
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
			table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
			table.SetCenterSeparator("|")
			table.Append(userData)

			table.Render()
		}
	}
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.Flags().BoolVarP(&all, "all", "a", false, "list all user parameters")
}
