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
	injectTaskData(testDB)
	tasks, err := repository.ReadTaskToday(testDB)
	if err != nil {
		applog.Error(err, "Error in repository.ReadTask test")
		t.Fatal()
	}
	for _, toPrint := range tasks {
		fmt.Printf("%v %v %v \n", toPrint.CreatedAt, toPrint.Description, toPrint.Status)
	}
}

func TestWriteTask(t *testing.T) {
	task := randomTaskData()
	err := repository.WriteTask(testDB, task)
	if err != nil {
		applog.Error(err, "Error in repository.WriteTask test")
		t.Fatal()
	}
	fmt.Printf("\n%v %v %v \n", task.CreatedAt, task.Description, task.Status)
}

func TestEditTask(t *testing.T) {
	task, _ := injectTaskData(testDB)
	fmt.Printf("%v %v %v \n", task.ID, task.Description, task.Status)
	repository.EditTask(testDB, &task, models.StatusDone)

	// en passant
	taskEdited, _ := repository.FindById(testDB, task.ID)
	applog.Debug(fmt.Sprintf("%v %v %v \n", task.ID, task.Description, task.Status))
	if taskEdited.Status != models.StatusDone {
		t.Error("CUNGPRET")
	}
}

func injectTaskData(testDB *gorm.DB) (models.Task, error) {
	task := randomTaskData()
	result := testDB.Create(&task)
	if result.Error != nil {
		applog.Error(result.Error, "Fail to injectTaskData")
		return models.Task{}, result.Error
	}
	return task, nil
}

func randomTaskData() models.Task {
	randId := uint16(rand.Uint32())
	return models.Task{
		ID:          randId,
		Description: fmt.Sprintf("Random Task %v", randId),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     time.Now().Add(time.Hour * 24),
		Status:      models.StatusNotStarted,
	}
}
