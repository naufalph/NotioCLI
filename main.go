package main

import (
	"tdlst/cmd"
	"tdlst/db"
)

func main() {
	db.ConnectMain()
	cmd.Execute()
}
