package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/doneill/er-cli/api"
	"github.com/doneill/er-cli/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// ----------------------------------------------
// patrols command
// ----------------------------------------------

var patrolsCmd = &cobra.Command{
	Use:   "patrols",
	Short: "Get patrols data",
	Long:  `Return patrol data including ID, serial number, state, title, and leader`,
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
	patrolsResponse, err := client.Patrols()
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

func formatPatrolData(patrol *api.Patrol) []string {
	leader := "N/A"
	if len(patrol.PatrolSegments) > 0 && patrol.PatrolSegments[0].Leader != nil {
		l := patrol.PatrolSegments[0].Leader
		leader = fmt.Sprintf("%s %s (%s)", l.FirstName, l.LastName, l.Username)
	}

	title := "N/A"
	if patrol.Title != nil {
		title = *patrol.Title
	}

	return []string{
		patrol.ID,
		fmt.Sprintf("%d", patrol.SerialNumber),
		patrol.State,
		title,
		leader,
	}
}

func configurePatrolsTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID",
		"Serial",
		"State",
		"Title",
		"Leader",
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
	rootCmd.AddCommand(patrolsCmd)
}
