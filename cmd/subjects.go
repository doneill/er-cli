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

var daysAgo int

// ----------------------------------------------
// subjects command
// ----------------------------------------------

var subjectsCmd = &cobra.Command{
	Use:   "subjects",
	Short: "Get updated subjects data",
	Long:  `Return subject data updated within specified number of days ago`,
	Run: func(cmd *cobra.Command, args []string) {
		subjects()
	},
}

// ----------------------------------------------
// functions
// ----------------------------------------------

func subjects() {
	updatedSince := time.Now().AddDate(0, 0, -daysAgo).UTC().Format("2006-01-02T15:04:05.000")
	client := api.ERClient(config.Sitename(), config.Token())

	subjectsResponse, err := client.Subjects(updatedSince)
	if err != nil {
		log.Fatalf("Error getting subjects: %v", err)
	}

	if subjectsResponse == nil || len(subjectsResponse.Data) == 0 {
		fmt.Println("No subjects found")
		return
	}

	table := configureSubjectsTable()
	for _, subject := range subjectsResponse.Data {
		table.Append(formatSubjectData(&subject))
	}

	table.Render()
}

func formatSubjectData(subject *api.Subject) []string {
	coordinates := "N/A"

	if len(subject.LastPosition.Geometry.Coordinates) >= 2 {
		coordinates = fmt.Sprintf("%.6f, %.6f",
			subject.LastPosition.Geometry.Coordinates[1],
			subject.LastPosition.Geometry.Coordinates[0],
		)
	}

	return []string{
		subject.Name,
		subject.ID,
		subject.SubjectType,
		subject.SubjectSubtype,
		subject.LastPositionDate,
		coordinates,
	}
}

func configureSubjectsTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"ID",
		"Type",
		"Subtype",
		"Last Position Date",
		"Last Position (lat, lon)",
	})
	table.SetBorders(tablewriter.Border{
		Left:   true,
		Top:    true,
		Right:  true,
		Bottom: true,
	})
	table.SetCenterSeparator("|")
	return table
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(subjectsCmd)
	subjectsCmd.Flags().IntVarP(&daysAgo, "updated-since", "u", 3, "Number of days ago to query updates from")
}
