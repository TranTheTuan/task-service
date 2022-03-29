package usecase

import (
	"context"

	"github.com/TranTheTuan/task-service/app/domain/dto"
	"github.com/TranTheTuan/task-service/app/domain/models"
	"github.com/TranTheTuan/task-service/app/domain/service"
)

type TaskUsecaseInterface interface {
	GetTasks(ctx context.Context, userID uint32, getTasksDto *dto.GetTasksDTO) ([]*models.Task, *models.Pagination, error)
	GetTaskByID(ctx context.Context, userID, taskID uint32) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTaskByID(ctx context.Context, userID, taskID uint32) error
}

type TaskUsecase struct {
	taskService service.TaskServiceInterface
}

func NewTaskUsecase(taskService service.TaskServiceInterface) *TaskUsecase {
	return &TaskUsecase{taskService}
}

func (u *TaskUsecase) GetTasks(ctx context.Context, userID uint32, getTasksDto *dto.GetTasksDTO) ([]*models.Task, *models.Pagination, error) {
	return u.taskService.GetTasks(ctx, userID, getTasksDto)
}

func (u *TaskUsecase) GetTaskByID(ctx context.Context, userID, taskID uint32) (*models.Task, error) {
	return u.taskService.GetTaskByID(ctx, userID, taskID)
}

func (u *TaskUsecase) CreateTask(ctx context.Context, task *models.Task) error {
	return u.taskService.CreateTask(ctx, task)
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, task *models.Task) error {
	return u.taskService.UpdateTask(ctx, task)
}

func (u *TaskUsecase) DeleteTaskByID(ctx context.Context, userID, taskID uint32) error {
	return u.taskService.DeleteTaskByID(ctx, userID, taskID)
}
