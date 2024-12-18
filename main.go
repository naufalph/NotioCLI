package main

import (
	"tdlst/cmd"
	"tdlst/db"
	"tdlst/pkg/applog"
)

func main() {
	db.ConnectMain()
	applog.InitLog()
	cmd.Execute()
}
