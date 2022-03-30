package repository

import (
	"context"

	"github.com/TranTheTuan/task-service/app/domain/dto"
	"github.com/TranTheTuan/task-service/app/domain/models"
)

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, userID uint32, getTasksDto *dto.GetTasksDTO) ([]*models.Task, *models.Pagination, error)
	GetTaskByID(ctx context.Context, userID, taskID uint32) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTaskByID(ctx context.Context, userID, taskID uint32) error
}
