/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "NotioCLI",
	Short: "NotioCLI is a lightweight, command-line to-do list application built with Go, designed for seamless integration with Notion’s database. It enables users to manage tasks directly from the terminal, providing a fast, distraction-free workflow for task organization",
	Long: `NotioCLI is a lightweight, command-line to-do list application built with Go, designed for seamless integration with Notion’s database. It enables users to manage tasks directly from the terminal, providing a fast, distraction-free workflow for task organization
	Key features include:

	- Create, Read, Update, Delete (CRUD) functionality for tasks in a Notion table.
	- Synchronization between local tasks and Notion’s database for online/offline flexibility.
	- Extensibility: Modular architecture built with cobra-cli, allowing easy addition of new features.
	- Fast and Reliable: Written in Go for performance and simplicity.
	- Secure: Uses environment variables for secure storage of API tokens.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.NotioCLI.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
