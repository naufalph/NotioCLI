package test

import (
	"fmt"
	"math/rand"
	"tdlst/internal/notion"
	m "tdlst/models"
	"tdlst/pkg/applog"
	"testing"
	"time"
)

func TestQueryNotion(t *testing.T) {
	tasks, err := notion.GetTaskList()
	if err != nil {
		applog.Error(err, "Error in notion.GetTaskList() test")
		t.Fatal()
	}
	for _, toPrint := range tasks {
		fmt.Printf("%v %v %v \n", toPrint.CreatedAt, toPrint.Description, toPrint.Status)
	}
}

func TestInsertNotion(t *testing.T) {
	randId := uint16(rand.Uint32())
	task := m.Task{
		ID:          randId,
		Description: fmt.Sprintf("New Task Test Insert %v", randId),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     time.Now().Add(time.Hour * 24),
		Status:      m.StatusNotStarted,
	}
	err := notion.AddTask(task)
	if err != nil {
		applog.Error(err, "Error in notion.AddTask test")
		t.Fatal()
	}
}

func TestUpdateNotion(t *testing.T) {
	task := m.Task{
		ID:          1,
		Description: "Aku adalah kelaparan itu sendiri",
		Status:      m.StatusInProgress,
		DueDate:     time.Now().Add(time.Hour * 24),
		NotionId:    "147dca9a-e8ac-8093-969a-fa5a605c9c89",
	}

	err := notion.UpdateTask(task)
	if err != nil {
		applog.Error(err, "Error in notion.AddTask test")
		t.Fatal()
	}
}

func TestSyncNotion(t *testing.T) {
	err := notion.SyncTask(testDB)
	if err != nil {
		applog.Error(err, "Error in notion.AddTask test")
		t.Fatal()
	}
}
