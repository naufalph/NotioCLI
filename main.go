/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"NotioCLI/cmd"
	"NotioCLI/db"
)

func main() {
	db.ConnectMain()
	cmd.Execute()
}
