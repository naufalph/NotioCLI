package internal

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"tdlst/internal/repository"
	m "tdlst/models"
	"tdlst/pkg/applog"
	"tdlst/pkg/utils"
	"tdlst/ui"
)

func UpdateTask(statusCode int8, ID uint16) error {
	var status m.TaskStatus
	if !slices.Contains(utils.AvailableStatusCodes, statusCode) {
		applog.Error(errors.New("input error"), "Unavailable Status Code")
		return errors.New("input error")
	} else {
		applog.Debug(fmt.Sprintf("TASDAS : %v", statusCode))
		status = utils.ParseStatusCode(statusCode)
	}

	if ID == 0 {
		applog.Error(errors.New("input error"), "Null ID")
		return errors.New("input error")
	}
	applog.Debug(fmt.Sprintf("ID to search :%d", ID))
	task, err := repository.FindById(nil, ID)
	if err != nil {
		applog.Error(err, "Wrong ID")
		return errors.New("input error")
	}
	return repository.EditTaskStatus(nil, task, status)
}

func DefUpdateExplanation() {
	ui.PrintLine(`
Invalid Input, to write the update command use task' ID to edit
Example:
tdlist update 199221
-> tdlist update [ID]

ps. task ID can be found on tdlst list command
	`)
}

func DefUpdateSuccessMsg(ID string, statusCode int8) {
	ui.PrintLine(fmt.Sprintf(`
Task %v is succesfully updated to %v status
`, ID, utils.ParseStatusCode(statusCode)))
}

func ParseIDFromArgs(ID string) uint16 {
	if ID == "" {
		applog.Error(errors.New("invalid input"), "Invalid Input")
		ui.PrintLine("Invalid Input")
	}

	IDuint, err := strconv.ParseUint(ID, 10, 32)
	if err != nil {
		applog.Error(err, "Invalid Input")
		ui.PrintLine("Invalid Input")
	}
	return uint16(IDuint)
}
