package internal

import (
	"math/rand/v2"
	"tdlst/internal/repository"
	m "tdlst/models"
	"tdlst/pkg/applog"
	"tdlst/ui"
	"time"
)

func WriteTask(task []string) {
	taskParsed, err := parseTask(task)
	if err != nil {
		applog.Error(err, "Error in parsing WriteTask test")
	}
	err = repository.WriteTask(nil, taskParsed)
	if err != nil {
		applog.Error(err, "Error in repository.WriteTask test")
	}
	ui.PrintLine("Task has been added !")
}

func DefWriteExplanation() {
	ui.PrintLine(`Invalid Input, to write the new command use description after the command
Example:
tdlist new "Harus bangun pagi" 2006-01-02
-> tdlist new [description] [task due date]
	`)
}

func parseTask(task []string) (m.Task, error) {
	randId := uint16(rand.Uint32())
	var dueDate time.Time
	var err error
	var dueDateParsed time.Time
	if len(task) < 2 {
		dueDate = time.Now().Add(24 * time.Hour)
	} else {
		dueDateParsed, err = time.Parse("2006-01-02", task[1])
		if err != nil {
			ui.PrintLine("Invalid Input!")
			dueDate = time.Now().Add(24 * time.Hour)
		} else {
			dueDate = dueDateParsed
		}
	}
	return m.Task{
		ID:          randId,
		Description: task[0],
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     dueDate,
		Status:      m.StatusNotStarted,
	}, err
}
