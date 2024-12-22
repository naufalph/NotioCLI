/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"tdlst/internal"
	"tdlst/pkg/applog"
	"tdlst/ui"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [task ID]",
	Short: "Update the status of an existing task",
	Long: `The "update" command allows you to modify the status of a task in your task manager.
	
	Provide the task ID as the first argument to identify which task you want to update. 
	You will then be prompted to select the new status for the task from the following options:
	
	0: Not Started
	1: In Progress
	2: Done
	3: Overdue
	4: Overdue Extra
	
	Example usage:
	  tdlst update 123
	This will allow you to update the status of the task with ID 123.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			internal.DefUpdateExplanation()
			applog.Error(errors.New("ID is required"), "Input ID is required")
		}
		IDuint16 := internal.ParseIDFromArgs(args[0])
		ui.PrintLine(`
Please choose your choosing status :
0 NotStarted
1 InProgress
2 Done
3 OverDue
4 OverDueExtra
		`)
		var status int
		fmt.Scanf("%d", &status)
		statusCode := int8(status)
		internal.UpdateTask(statusCode, IDuint16)
		internal.DefUpdateSuccessMsg(args[0], statusCode)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
