/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"tdlst/internal"
	"tdlst/pkg/applog"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Command to add your new task",
	Long: `The "new" command allows you to create a new task in your Notio workspace.
	You can specify task details for description, and due date directly from the command line.
	
	Examples:
	  notiocli new "Write blog post" 2024-12-17 --priority high
	  notiocli new "Weekly team meeting" --desc "Prepare the agenda and notes" --priority medium`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			applog.Error(nil, "Description is required")
		}
		internal.WriteTask(args)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
