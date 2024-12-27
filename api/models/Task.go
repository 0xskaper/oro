package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Title      string    `gorm:"size:20;not null" json:"title"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsComplete bool      `gorm:"default:false" json:"is_complete"`
	TaskListID uint32    `json:"task_list_id"`
}

func (t *Task) Prepare() {
	t.ID = 0
	t.IsComplete = false
	t.UpdatedAt = time.Now()
	t.CreatedAt = time.Now()
	t.Title = html.EscapeString(strings.TrimSpace(t.Title))
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("required title")
	}
	if t.TaskListID < 1 {
		return errors.New("required task list ID")
	}
	return nil
}

func (t *Task) CreateNewTask(db *gorm.DB) (*Task, error) {
	tx := db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return &Task{}, tx.Error
	}

	var taskList TaskList
	err := tx.Model(&TaskList{}).Where("id = ?", t.TaskListID).First(&taskList).Error
	if err != nil {
		tx.Rollback()
		return &Task{}, errors.New("task list not found.")
	}

	err = tx.Model(&TaskList{}).Where("id = ?", t.TaskListID).UpdateColumn("total_tasks + ?", 1).Error
	err = tx.Model(&Task{}).Create(&t).Error

	if err != nil {
		tx.Rollback()
		return &Task{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return &Task{}, err
	}

	return t, nil
}

func UpdateTaskCompletionStatus(db *gorm.DB, taskID uint32, isTaskComplete bool) error {
	tx := db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	err := tx.Model(&Task{}).Where("id = ?", taskID).Update("is_complete", isTaskComplete).Error

	if err != nil {
		tx.Rollback()
		return tx.Error
	}

	return tx.Commit().Error
}
