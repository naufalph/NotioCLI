package repository

import (
	"tdlst/db"
	"tdlst/models"
	"tdlst/pkg/applog"
	"time"

	"gorm.io/gorm"
)

func ReadTaskToday(dbTest *gorm.DB) ([]models.Task, error) {
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

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
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
