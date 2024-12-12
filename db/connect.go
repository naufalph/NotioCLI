package db

import (
	"fmt"
	"tdlst/config"
	"tdlst/pkg/applog"
	"tdlst/pkg/utils"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(envFileLoc string) (db *gorm.DB, err error) {
	applog.Debug("Starting to connect ...")
	dsn := getDsn(envFileLoc)
	applog.Debug("Connecting to " + dsn)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func ConnectMain() {
	db, err := gorm.Open(mysql.Open(getDsn(config.RealEnv())), &gorm.Config{})
	if err != nil {
		applog.Error(err, utils.DBConnectionError)
	}
	DB = db
}

func getDsn(envFileLoc string) string {

	envMap, err := godotenv.Read(envFileLoc)

	if err != nil {
		applog.Error(err, err.Error())
	}

	user := envMap["MYSQL_USER"]
	pass := envMap["MYSQL_PASSWORD"]
	host := envMap["MYSQL_DB_HOST"]
	dbname := envMap["MYSQL_DB_NAME"]
	protocol := envMap["MYSQL_PROTOCOL"]

	return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true", user, pass, protocol, host, dbname)

}
