package test

import (
	"fmt"
	"math/rand"
	"os"
	"tdlst/config"
	"tdlst/db"
	"tdlst/internal/repository"
	"tdlst/models"
	"tdlst/pkg/applog"
	"tdlst/pkg/utils"
	"testing"
	"time"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = db.Connect(config.DevEnv())
	if err != nil {
		applog.Error(err, utils.DBConnectionError)
		os.Exit(1)
	}
	testDB.AutoMigrate(&models.Task{})
	code := m.Run()

	sqlDb, _ := testDB.DB()
	sqlDb.Close()

	os.Exit(code)
}
func TestReadTask(t *testing.T) {
	_ = injectTaskData(testDB)
	tasks, err := repository.ReadTaskToday(testDB)
	if err != nil {
		applog.Error(err, "Error in repository.ReadTask test")
		t.Fatal()
	}
	for _, toPrint := range tasks {
		fmt.Printf("%v %v %v \n", toPrint.CreatedAt, toPrint.Description, toPrint.Status)
	}
}

func injectTaskData(testDB *gorm.DB) error {
	task := randomTaskData()
	result := testDB.Create(&task)
	if result.Error != nil {
		applog.Error(result.Error, "Fail to injectTaskData")
		return result.Error
	}
	return nil
}

func randomTaskData() models.Task {
	randId := uint16(rand.Uint32())
	randDesc := fmt.Sprintf("Random Task %v", randId)
	status := utils.StatusNotStarted
	return models.Task{
		ID:          randId,
		Description: randDesc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     time.Now().Add(time.Hour * 24),
		Status:      status,
	}
}
