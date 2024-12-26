package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/doneill/er-cli/api"
	"github.com/doneill/er-cli/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	daysAgo   int
	tracks    bool
	subjectID string
	export    bool
)

// ----------------------------------------------
// subjects command
// ----------------------------------------------

var subjectsCmd = &cobra.Command{
	Use:   "subjects",
	Short: "Get subjects data and tracks",
	Long:  `Return subject data and tracks updated within specified number of days ago`,
	Run: func(cmd *cobra.Command, args []string) {
		subjects()
	},
}

// ----------------------------------------------
// functions
// ----------------------------------------------

func subjects() {
	client := api.ERClient(config.Sitename(), config.Token())

	switch {
	case tracks:
		if subjectID == "" {
			log.Fatal("Subject ID is required when using tracks flag")
		}
		handleTracks(client)
	default:
		handleSubjects(client)
	}
}

func handleTracks(client *api.Client) {
	tracksResponse, err := client.SubjectTracks(subjectID, daysAgo)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Fatalf("Subject ID %s not found. Use 'er subjects' command to list available subjects", subjectID)
		}
		log.Fatalf("Error getting tracks: %v", err)
	}

	if tracksResponse == nil || len(tracksResponse.Data.Features) == 0 {
		fmt.Println("No tracks found")
		return
	}

	// Create a GeoJSON structure
	geoJSON := struct {
		Type     string        `json:"type"`
		Features []api.Feature `json:"features"`
	}{
		Type:     tracksResponse.Data.Type,
		Features: tracksResponse.Data.Features,
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(geoJSON, "", "    ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}

	if export {
		filename := fmt.Sprintf("%s.geojson", subjectID)
		if err := os.WriteFile(filename, jsonData, 0644); err != nil {
			log.Fatalf("Error writing file: %v", err)
		}
		fmt.Printf("Successfully exported tracks to %s\n", filename)
	} else {
		fmt.Println(string(jsonData))
	}
}

func handleSubjects(client *api.Client) {
	// Calculate the date from days ago
	updatedSince := time.Now().AddDate(0, 0, -daysAgo).UTC().Format("2006-01-02T15:04:05.000")

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
	subjectsCmd.Flags().BoolVarP(&tracks, "tracks", "t", false, "Get tracks for a subject")
	subjectsCmd.Flags().StringVarP(&subjectID, "subject-id", "s", "", "Subject ID for tracks query")
	subjectsCmd.Flags().BoolVarP(&export, "export", "e", false, "Export tracks to GeoJSON file")
}
