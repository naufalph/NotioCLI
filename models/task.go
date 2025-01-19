package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          uint16 `gorm:"primary key"`
	Description string `gorm:"size:255;not null"`
	DueDate     time.Time
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	Status      TaskStatus `gorm:"default:'Not Started'"`
	NotionId    string     `gorm:"size:255;default:''"`
}

type TaskStatus string

const (
	StatusNotStarted   TaskStatus = "Not Started"
	StatusInProgress   TaskStatus = "In Progress"
	StatusDone         TaskStatus = "Done"
	StatusOverDue      TaskStatus = "Overdue"
	StatusOverDueExtra TaskStatus = "Overdue Dude"
)

var AvailableStatuses = []TaskStatus{
	StatusNotStarted,
	StatusInProgress,
	StatusDone,
	StatusOverDue,
	StatusOverDueExtra,
}
