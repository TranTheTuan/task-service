package models

import (
	"github.com/jinzhu/gorm"

	pbTasks "github.com/TranTheTuan/pbtypes/build/go/tasks"
)

type Task struct {
	gorm.Model
	UserID      uint32 `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (t *Task) ToProto() *pbTasks.Task {
	return &pbTasks.Task{
		Id:          uint32(t.ID),
		Name:        t.Name,
		Description: t.Description,
	}
}
