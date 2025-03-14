package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/doneill/er-cli/api"
	"github.com/doneill/er-cli/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	all      bool
	profiles bool
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
	er := api.ERClient(config.Sitename(), config.Token())
	userResponse, err := er.User()
	if err != nil {
		log.Fatalf("Error authentication: %v", err)
	}
	if userResponse == nil {
		log.Fatalf("No user response received")
	}

	switch {
	case profiles:
		handleProfiles(userResponse.Data.ID)
	case all:
		fmt.Println(formatAllUserData(userResponse))
	default:
		table := configureTable()
		table.Append(formatTableData(userResponse))
		table.Render()
	}

	if err := updateConfig(userResponse.Data.ID); err != nil {
		log.Printf("Warning: failed to update config: %v", err)
	}
}

func handleProfiles(userID string) {
	er := api.ERClient(config.Sitename(), config.Token())
	profilesResponse, err := er.UserProfiles(userID)
	if err != nil {
		log.Fatalf("Error getting profiles: %v", err)
	}

	if profilesResponse == nil || len(profilesResponse.Data) == 0 {
		fmt.Println("No profiles found")
		return
	}

	switch {
	case all:
		for i, profile := range profilesResponse.Data {
			if i > 0 {
				fmt.Println("\n---")
			}
			fmt.Println(formatAllProfileData(&profile))
		}
	default:
		table := configureTable()
		for _, profile := range profilesResponse.Data {
			table.Append(formatTableData(&api.UserResponse{Data: profile}))
		}
		table.Render()
	}
}

func formatAllUserData(data *api.UserResponse) string {
	return fmt.Sprintf("username: %s\nemail: %s\nfirst name: %s\n"+
		"last name: %s\nrole: %s\nis staff: %t\nis superuser: %t\n"+
		"date joined: %s\nid: %s\nisactive: %t\nlast login: %s\n"+
		"pin: %s\nsubject id: %s",
		data.Data.Username,
		data.Data.Email,
		data.Data.FirstName,
		data.Data.LastName,
		data.Data.Role,
		data.Data.IsStaff,
		data.Data.IsSuperUser,
		data.Data.DateJoined,
		data.Data.ID,
		data.Data.IsActive,
		data.Data.LastLogin,
		data.Data.Pin,
		data.Data.Subject.ID)
}

func formatAllProfileData(data *api.UserData) string {
	return fmt.Sprintf("username: %s\nemail: %s\nfirst name: %s\n"+
		"last name: %s\nrole: %s\nis staff: %t\nis superuser: %t\n"+
		"date joined: %s\nid: %s\nisactive: %t\nlast login: %s\n"+
		"pin: %s\nsubject id: %s\naccepted eula: %t",
		data.Username,
		data.Email,
		data.FirstName,
		data.LastName,
		data.Role,
		data.IsStaff,
		data.IsSuperUser,
		data.DateJoined,
		data.ID,
		data.IsActive,
		data.LastLogin,
		data.Pin,
		data.Subject.ID,
		data.AcceptedEula)
}

func formatTableData(data *api.UserResponse) []string {
	return []string{
		data.Data.Username,
		data.Data.Email,
		data.Data.FirstName,
		data.Data.LastName,
		data.Data.ID,
		data.Data.Pin,
		data.Data.Subject.ID,
	}
}

func configureTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Username", "Email", "First Name", "Last Name", "ID", "Pin", "Subject ID"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetCenterSeparator("|")
	return table
}

func updateConfig(userID string) error {
	if viper.IsSet("remote_id") {
		return nil
	}

	viper.Set("remote_id", userID)
	if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.Flags().BoolVarP(&all, "all", "a", false, "list all user parameters")
	userCmd.Flags().BoolVarP(&profiles, "profiles", "p", false, "list all user profiles")
}
