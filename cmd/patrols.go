package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/doneill/er-cli/api"
	"github.com/doneill/er-cli/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var days int

// ----------------------------------------------
// patrols command
// ----------------------------------------------

var patrolsCmd = &cobra.Command{
	Use:   "patrols",
	Short: "Get patrols data",
	Long:  `Return patrol data including serial number, state, ID, location, and time information`,
	Run: func(cmd *cobra.Command, args []string) {
		patrols()
	},
}

// ----------------------------------------------
// functions
// ----------------------------------------------

func patrols() {
	client := api.ERClient(config.Sitename(), config.Token())
	handlePatrols(client)
}

func handlePatrols(client *api.Client) {
	patrolsResponse, err := client.Patrols(days)
	if err != nil {
		log.Fatalf("Error getting patrols: %v", err)
	}

	if patrolsResponse == nil || len(patrolsResponse.Data.Results) == 0 {
		fmt.Println("No patrols found")
		return
	}

	table := configurePatrolsTable()
	for _, patrol := range patrolsResponse.Data.Results {
		table.Append(formatPatrolData(&patrol))
	}
	table.Render()
}

func formatTime(timeStr *string) string {
	if timeStr == nil {
		return "N/A"
	}

	t, err := time.Parse(time.RFC3339, *timeStr)
	if err != nil {
		return "Invalid Time"
	}

	return t.Format("02 Jan 15:04")
}

func formatPatrolData(patrol *api.Patrol) []string {
	leader := "N/A"
	location := "N/A"
	startTime := "N/A"
	endTime := "N/A"

	if len(patrol.PatrolSegments) > 0 {
		segment := patrol.PatrolSegments[0]

		if segment.Leader != nil {
			l := segment.Leader
			leader = l.Name
		}

		if segment.StartLocation != nil {
			location = fmt.Sprintf("%.6f, %.6f",
				segment.StartLocation.Latitude,
				segment.StartLocation.Longitude)
		}

		startTime = formatTime(segment.TimeRange.StartTime)
		endTime = formatTime(segment.TimeRange.EndTime)
	}

	title := "N/A"
	if patrol.Title != nil {
		title = *patrol.Title
	}

	return []string{
		fmt.Sprintf("%d", patrol.SerialNumber),
		patrol.State,
		patrol.ID,
		title,
		leader,
		location,
		startTime,
		endTime,
	}
}

func configurePatrolsTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Serial",
		"State",
		"ID",
		"Title",
		"Leader",
		"Start Location",
		"Start Time",
		"End Time",
	})
	table.SetBorders(tablewriter.Border{
		Left:   true,
		Top:    true,
		Right:  true,
		Bottom: true,
	})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	return table
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(patrolsCmd)
	patrolsCmd.Flags().IntVarP(&days, "days", "d", 0, "Number of days to fetch patrols for (defaults to all)")
}
