package models

import "time"

type Filter struct {
	fromDate time.Time
	toDate   time.Time
}
