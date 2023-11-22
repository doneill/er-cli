package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ----------------------------------------------
// const var
// ----------------------------------------------

const PROGRAM_NAME string = "er"
const CONFIG_TYPE string = ".toml"
const CONFIG_PATH string = "/Users/dano/.config/er/"

// ----------------------------------------------
// command
// ----------------------------------------------

var rootCmd = &cobra.Command{
	Use:     "er",
	Short:   "EarthRanger CLI",
	Long:    `Work with EarthRanger platform from command line`,
	Version: "0.1.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}


func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName(PROGRAM_NAME)
	viper.SetConfigType(CONFIG_TYPE)
	viper.AddConfigPath(CONFIG_PATH)
}
