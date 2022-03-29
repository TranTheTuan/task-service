package service

import (
	"context"

	"github.com/TranTheTuan/task-service/app/domain/dto"
	"github.com/TranTheTuan/task-service/app/domain/models"
	"github.com/TranTheTuan/task-service/app/domain/repository"
)

type TaskServiceInterface interface {
	GetTasks(ctx context.Context, userID uint32, getTasksDto *dto.GetTasksDTO) ([]*models.Task, *models.Pagination, error)
	GetTaskByID(ctx context.Context, userID, taskID uint32) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTaskByID(ctx context.Context, userID, taskID uint32) error
}

type TaskService struct {
	taskRepo repository.TaskRepositoryInterface
}

func NewTaskService(taskRepo repository.TaskRepositoryInterface) *TaskService {
	return &TaskService{taskRepo}
}

func (s *TaskService) GetTasks(ctx context.Context, userID uint32, getTasksDto *dto.GetTasksDTO) ([]*models.Task, *models.Pagination, error) {
	return s.taskRepo.GetTasks(ctx, userID, getTasksDto)
}

func (s *TaskService) GetTaskByID(ctx context.Context, userID, taskID uint32) (*models.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, userID, taskID)
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	return s.taskRepo.UpdateTask(ctx, task)
}

func (s *TaskService) DeleteTaskByID(ctx context.Context, userID, taskID uint32) error {
	return s.taskRepo.DeleteTaskByID(ctx, userID, taskID)
}
