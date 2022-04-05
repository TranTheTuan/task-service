package models

import (
	"github.com/jinzhu/gorm"

	"github.com/TranTheTuan/pbtypes/build/go/core"
)

type Task struct {
	gorm.Model
	UserID      uint32 `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (t *Task) ToProto() *core.Task {
	return &core.Task{
		Id:          uint32(t.ID),
		Name:        t.Name,
		Description: t.Description,
	}
}
