package test

import (
	"tdlst/config"
	"tdlst/db"
	"tdlst/pkg/applog"
	"testing"
)

func TestTestDBConnection(t *testing.T) {
	applog.Debug("Using env : config.DevEnv()")
	db, err := db.Connect(config.DevEnv())
	if err != nil {
		applog.Error(err, "Database connection failed")
		t.Fatal()
	}
	sqlDb, err := db.DB()
	if err != nil {
		applog.Error(err, "Get database failed")
		t.Fatal()
	}
	defer sqlDb.Close()
	err = sqlDb.Ping()
	if err != nil {
		applog.Error(err, "Ping db error")
		t.Fatal()
		return
	}
	applog.Info(nil, "DB Connection Succeeded!")
}
