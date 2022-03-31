package repository

import (
	"context"
	"math"

	"github.com/jinzhu/gorm"

	"github.com/TranTheTuan/task-service/app/domain/dto"
	"github.com/TranTheTuan/task-service/app/domain/models"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (repo *TaskRepository) GetTasks(ctx context.Context, userID uint32, getTasksDto *dto.GetTasksDTO) ([]*models.Task, *models.Pagination, error) {
	var tasks []*models.Task
	limit := getTasksDto.Limit
	page := getTasksDto.Page
	keyword := getTasksDto.Keyword
	sort := getTasksDto.Sort
	if limit <= 0 {
		limit = 1
	}
	if page <= 0 {
		page = 1
	}
	if sort == "" {
		sort = "id desc"
	}
	offset := (page - 1) * limit
	db := repo.db.Model(&models.Task{}).Where("user_id = ?", userID)
	if keyword != "" {
		keyword = "%" + keyword + "%"
		db = db.Where("name like ? OR description like ?", keyword, keyword)
	}

	var total uint32
	if err := db.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	db = db.Offset(offset).Limit(limit).Order(sort).Find(&tasks)
	err := db.Error
	if err != nil {
		return nil, nil, err
	}

	pagination := &models.Pagination{
		Limit:      limit,
		Page:       page,
		Total:      uint32(total),
		TotalPages: uint32(math.Ceil(float64(total) / float64(limit))),
	}
	return tasks, pagination, nil
}

func (repo *TaskRepository) GetTaskByID(ctx context.Context, userID, taskID uint32) (*models.Task, error) {
	task := &models.Task{}
	if err := repo.db.Model(&models.Task{}).Where("user_id = ? AND id = ?", userID, taskID).First(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (repo *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	return repo.db.Model(&models.Task{}).Create(task).Error
}

func (repo *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	return repo.db.Model(&models.Task{}).Where("user_id = ?", task.UserID).Update(task).Error
}

func (repo *TaskRepository) DeleteTaskByID(ctx context.Context, userID, taskID uint32) error {
	return repo.db.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskID, userID).Delete(&models.Task{}).Error
}
