package cmd

import (
	"fmt"

	"github.com/doneill/er-cli-go/data"
	"github.com/spf13/cobra"
)

var displayTables bool

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
// funtions
// ----------------------------------------------

func open(file string) {
	db, err := data.DbConnect(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch {
	case displayTables:
		tables, err := data.GetTables(*db)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, tableName := range tables {
			fmt.Println(tableName)
		}
	default:
		message := fmt.Sprintf("%s successfully opened", file)
		fmt.Println(message)
	}
}

// ----------------------------------------------
// initialize
// ----------------------------------------------

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().BoolVarP(&displayTables, "tables", "t", false, "Display all database tables")
}
