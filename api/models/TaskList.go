package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type TaskList struct {
	ID                 uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name               string `gorm:"size:10;not null" json:"name"`
	TotalCompleteTasks int64  `json:"total_complete_task"`
	TotalTasks         int64  `json:"total_tasks"`
}

func (tl *TaskList) Prepare() {
	tl.Name = html.EscapeString(strings.TrimSpace(tl.Name))
	tl.TotalCompleteTasks = 0
	tl.TotalTasks = 0
	tl.ID = 0
}

func (tl *TaskList) Validate() error {
	if tl.Name == "" {
		return errors.New("required tasklist name")
	}
	return nil
}

func (tl *TaskList) CreateNewTaskList(db *gorm.DB) (*TaskList, error) {
	var err error
	err = db.Debug().Model(&TaskList{}).Create(&tl).Error
	if err != nil {
		return &TaskList{}, err
	}
	return tl, nil
}

func (tl *TaskList) DeleteTaskList(db *gorm.DB, id uint32) (int64, error) {
	db = db.Debug().Model(&TaskList{}).Where("id = ?", id).Take(&TaskList{}).Delete(&TaskList{})
	if db.Error != nil {
		return 0, errors.New("No Task list found.")
	}
	return db.RowsAffected, nil
}
