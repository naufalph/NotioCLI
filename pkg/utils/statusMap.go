package utils

import (
	m "tdlst/models"
	"tdlst/pkg/applog"
)

const (
	StatusNotStartedCode   int8 = 0
	StatusInProgressCode   int8 = 1
	StatusDoneCode         int8 = 2
	StatusOverDueCode      int8 = 3
	StatusOverDueExtraCode int8 = 4
)

var AvailableStatusCodes = []int8{
	StatusNotStartedCode,
	StatusInProgressCode,
	StatusDoneCode,
	StatusOverDueCode,
	StatusOverDueExtraCode,
}

func ParseStatusCode(statusCode int8) m.TaskStatus {
	switch statusCode {
	case StatusNotStartedCode:
		return m.StatusNotStarted
	case StatusDoneCode:
		return m.StatusDone
	case StatusInProgressCode:
		return m.StatusInProgress
	case StatusOverDueCode:
		return m.StatusOverDue
	case StatusOverDueExtraCode:
		return m.StatusOverDueExtra
	}
	applog.Debug("Status Code unrecognized")
	return m.StatusNotStarted
}

func ParseStatusString(status string) m.TaskStatus {
	switch status {
	case string(m.StatusNotStarted):
		return m.StatusNotStarted
	case string(m.StatusDone):
		return m.StatusDone
	case string(m.StatusInProgress):
		return m.StatusInProgress
	case string(m.StatusOverDue):
		return m.StatusOverDue
	case string(m.StatusOverDueExtra):
		return m.StatusOverDueExtra
	}
	applog.Debug("Status Code unrecognized")
	return m.StatusNotStarted
}
