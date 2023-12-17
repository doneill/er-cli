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
	var count int64

	db, err := data.DbConnect(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch {
	case dbUser:
		var user []data.Accounts_User
		db.First(&user)
		fmt.Println(user[0].Username)
	case displayTables:
		tables, err := data.GetTables(*db)
		if err != nil {
			fmt.Println(err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Count"})

		for _, tableName := range tables {
			db.Table(tableName).Count(&count)
			table.Append([]string{tableName, fmt.Sprintf("%d", count)})
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
}
