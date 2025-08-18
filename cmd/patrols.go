package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/doneill/er-cli/api"
	"github.com/doneill/er-cli/config"
	"github.com/doneill/er-cli/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	days   int
	status string
	track  string
)

var validStatuses = map[string]bool{
	"active":    true,
	"done":      true,
	"cancelled": true,
}

// ----------------------------------------------
// patrols command
// ----------------------------------------------

var patrolsCmd = &cobra.Command{
	Use:   "patrols",
	Short: "Get patrols data",
	Long:  `Return patrol data including serial number, state, ID, location, and time information`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if status != "" {
			if !validStatuses[status] {
				return fmt.Errorf("invalid status value: %s\nValid status values are: active, done, cancelled", status)
			}
		}
		if days > 30 {
			return fmt.Errorf("days value cannot exceed 30 (got %d)", days)
		}
		if days < 0 {
			return fmt.Errorf("days value must be positive (got %d)", days)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if track != "" {
			patrolTrack()
		} else {
			patrols()
		}
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
	patrolsResponse, err := client.Patrols(days, status)
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
	startLocation := "N/A"
	endLocation := "N/A"
	startTime := "N/A"
	endTime := "N/A"
	segmentID := "N/A"

	if len(patrol.PatrolSegments) > 0 {
		segment := patrol.PatrolSegments[0]
		segmentID = segment.ID

		leader = segment.Leader.Name

		if segment.StartLocation != nil {
			startLocation = fmt.Sprintf("%.6f, %.6f",
				segment.StartLocation.Latitude,
				segment.StartLocation.Longitude)
		}

		if segment.EndLocation != nil {
			endLocation = fmt.Sprintf("%.6f, %.6f",
				segment.EndLocation.Latitude,
				segment.EndLocation.Longitude)
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
		startLocation,
		endLocation,
		startTime,
		endTime,
		segmentID,
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
		"End Location",
		"Start Time",
		"End Time",
		"Segment ID",
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

func patrolTrack() {
	client := api.ERClient(config.Sitename(), config.Token())
	handlePatrolTrack(client)
}

func handlePatrolTrack(client *api.Client) {
	// First get the patrol details
	patrolResponse, err := client.PatrolByID(track)
	if err != nil {
		log.Fatalf("Error getting patrol: %v", err)
	}

	if len(patrolResponse.Data.PatrolSegments) == 0 {
		fmt.Println("No patrol segments found for this patrol")
		return
	}

	// Get the first segment
	segment := patrolResponse.Data.PatrolSegments[0]
	
	if segment.TimeRange.StartTime == nil {
		fmt.Println("No start time found for patrol segment")
		return
	}

	// Determine end time - use current time if patrol is active (no end time)
	endTime := ""
	if segment.TimeRange.EndTime != nil {
		endTime = *segment.TimeRange.EndTime
	} else {
		endTime = time.Now().UTC().Format("2006-01-02T15:04:05Z")
		fmt.Printf("Patrol is active, using current time as end: %s\n", endTime)
	}

	// Get tracks for the patrol leader
	tracksResponse, err := client.PatrolTracks(segment.Leader.ID, *segment.TimeRange.StartTime, endTime)
	if err != nil {
		log.Fatalf("Error getting patrol tracks: %v", err)
	}

	if len(tracksResponse.Data.Features) == 0 {
		fmt.Println("No tracks found for this patrol")
		return
	}

	// Export to GeoJSON file
	filename := fmt.Sprintf("patrol_%d_tracks.geojson", patrolResponse.Data.SerialNumber)
	if err := utils.ExportToFile(tracksResponse.Data, filename); err != nil {
		log.Fatalf("Error exporting GeoJSON: %v", err)
	}

	fmt.Printf("Patrol tracks exported to %s\n", filename)
	fmt.Printf("Patrol: %s (Serial: %d)\n", patrolResponse.Data.Title, patrolResponse.Data.SerialNumber)
	fmt.Printf("Leader: %s\n", segment.Leader.Name)
	fmt.Printf("Start: %s\n", formatTime(segment.TimeRange.StartTime))
	fmt.Printf("End: %s\n", formatTime(segment.TimeRange.EndTime))
	fmt.Printf("Tracks: %d features\n", len(tracksResponse.Data.Features))
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(patrolsCmd)
	patrolsCmd.Flags().IntVarP(&days, "days", "d", 7, "Number of days to fetch patrols for")
	patrolsCmd.Flags().StringVarP(&status, "status", "s", "", "Patrol status (active, done, or cancelled)")
	patrolsCmd.Flags().StringVarP(&track, "track", "t", "", "Get tracks for a specific patrol ID")
}
