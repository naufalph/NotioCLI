package internal

import (
	"tdlst/internal/repository"
	m "tdlst/models"
	"tdlst/pkg/applog"
	"tdlst/pkg/utils"
	"tdlst/ui"
)

func List(listType string, filter m.Filter) {
	tasks, err := repository.ReadTaskToday(nil)
	if err != nil {
		applog.Error(err, "Error in repository.ReadTask test")
	}
	ui.PrintListNew(utils.ListAll, tasks)
}
