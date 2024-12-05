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
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	Status      string    `gorm:"default:false"`
}
