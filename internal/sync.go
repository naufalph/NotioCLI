package internal

import (
	"tdlst/internal/notion"
	"tdlst/pkg/applog"
	"tdlst/ui"
)

func Sync() error {
	err := notion.SyncTask(nil)
	if err != nil {
		applog.Error(err, "Error sync tasks")
	}
	ui.PrintLine("Task to Notion Synced")
	return err
}
