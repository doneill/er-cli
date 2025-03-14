package cmd

import (
	"fmt"
	"os"

	"github.com/doneill/er-cli/data"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var dbUser bool
var displayTables bool
var eventPendingSyncs bool
var eventDrafts bool

// ----------------------------------------------
// open command
// ----------------------------------------------

var openCmd = &cobra.Command{
	Use:   "open [sqlite db file]",
	Short: "Open a SQLite database file",
	Long:  `This tool is intended to be used specficially with EarthRanger mobile databases`,
	Run: func(cmd *cobra.Command, args []string) {
		var filename = args[0]
		open(filename)
	},
}

// ----------------------------------------------
// functions
// ----------------------------------------------

func open(file string) {
	data.OpenDB(file)

	switch {
	case dbUser:
		user := data.SelectUser()
		fmt.Println(user.Username)
	case displayTables:
		tables, err := data.GetTables()
		if err != nil {
			fmt.Println(err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Count"})

		for _, tableName := range tables {
			count := data.GetTableRowCount(tableName)
			table.Append([]string{tableName, fmt.Sprintf("%d", count)})
		}

		table.Render()
	case eventPendingSyncs:
		var users []string

		events := data.SelectPendingSyncEvents()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "User", "Title", "Values", "Patrol Segment ID", "Created At", "Updated At"})

		for _, event := range events {
			if event.ProfileID != 0 {
				profile := data.SelectUserProfileById(event.ProfileID)
				users = append(users, profile.Username)
			} else {
				user := data.SelectUser()
				users = append(users, user.Username)
			}
			table.Append([]string{fmt.Sprintf("%d", event.ID), users[len(users)-1], event.Title, event.Values, event.PatrolSegmentID, event.CreatedAt, event.UpdatedAt})
		}
		table.Render()
	case eventDrafts:
		var users []string

		events := data.SelectDrafts()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "User", "Title", "Created At", "Updated At"})

		for _, event := range events {
			if event.ProfileID != 0 {
				profile := data.SelectUserProfileById(event.ProfileID)
				users = append(users, profile.Username)
			} else {
				user := data.SelectUser()
				users = append(users, user.Username)
			}
			table.Append([]string{fmt.Sprintf("%d", event.ID), users[len(users)-1], event.Title, event.CreatedAt, event.UpdatedAt})
		}
		table.Render()
	default:
		message := fmt.Sprintf("%s successfully opened!", file)
		fmt.Println(message)
	}
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().BoolVarP(&dbUser, "user", "u", false, "Display database account user")
	openCmd.Flags().BoolVarP(&displayTables, "tables", "t", false, "Display all database tables")
	openCmd.Flags().BoolVarP(&eventPendingSyncs, "event-pending-syncs", "e", false, "Display all pending sync events")
	openCmd.Flags().BoolVarP(&eventDrafts, "event-drafts", "d", false, "Display all draft events")
}
