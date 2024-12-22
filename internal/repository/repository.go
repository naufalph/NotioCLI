package repository

import (
	"errors"
	"tdlst/db"
	"tdlst/models"
	m "tdlst/models"
	"tdlst/pkg/applog"
	"time"

	"gorm.io/gorm"
)

func ReadTaskToday(dbTest *gorm.DB) ([]m.Task, error) {
	dbUse := dbTest
	if dbTest == nil {
		dbUse = db.DB
	}
	startTime := time.Now().Truncate(24 * time.Hour)
	endTime := startTime.Add(24 * time.Hour)
	rows, err := dbUse.Raw(
		`SELECT id, created_at, updated_at, description, due_date, status   status FROM tasks 
		WHERE created_at BETWEEN ? AND ?`,
		startTime, endTime).Rows()
	if err != nil {
		applog.Error(err, "Error on gathering Task")
		return nil, err
	}
	defer rows.Close()

	var tasks []m.Task
	for rows.Next() {
		var task m.Task
		err := rows.Scan(
			&task.ID,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Description,
			&task.DueDate,
			&task.Status,
		)
		if err != nil {
			applog.Error(err, "Error on gathering Task")
			return nil, err
		} else {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func WriteTask(dbTest *gorm.DB, task m.Task) error {
	dbUse := dbTest
	if dbTest == nil {
		dbUse = db.DB
	}
	result := dbUse.Create(&task)
	return result.Error
}

func EditTask(dbTest *gorm.DB, task *m.Task, status m.TaskStatus) error {
	dbUse := dbTest
	if dbTest == nil {
		dbUse = db.DB
	}
	//test if exist
	if err := dbUse.First(&task).Error; err != nil {
		return err
	}
	task.Status = status
	return dbUse.Save(&task).Error
}

func FindById(dbTest *gorm.DB, ID uint16) (*m.Task, error) {
	dbUse := dbTest
	if dbTest == nil {
		dbUse = db.DB
	}
	var task models.Task
	result := dbUse.First(&task, ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("task not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}
