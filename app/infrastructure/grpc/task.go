package grpc

import (
	"context"

	"github.com/sirupsen/logrus"

	pbTasks "github.com/TranTheTuan/pbtypes/build/go/tasks"
	"github.com/TranTheTuan/pbtypes/build/go/core"
	"github.com/TranTheTuan/task-service/app/domain/dto"
	"github.com/TranTheTuan/task-service/app/domain/models"
	"github.com/TranTheTuan/task-service/app/domain/usecase"
	"github.com/TranTheTuan/task-service/app/infrastructure/util"
	pbTypes "google.golang.org/protobuf/types/known/emptypb"
)

type TaskServiceServer struct {
	taskUsecase usecase.TaskUsecaseInterface

	pbTasks.UnimplementedTaskCreateServiceServer
	pbTasks.UnimplementedTaskDeleteServiceServer
	pbTasks.UnimplementedTaskGetAllServiceServer
	pbTasks.UnimplementedTaskGetByIDServiceServer
	pbTasks.UnimplementedTaskUpdateServiceServer
}

func NewTaskServiceServer(taskUsecase usecase.TaskUsecaseInterface) *TaskServiceServer {
	return &TaskServiceServer{taskUsecase: taskUsecase}
}

func (t *TaskServiceServer) GetTasks(ctx context.Context, in *pbTasks.GetTasksRequest) (*pbTasks.GetTasksResponse, error) {
	userID, _ := util.GetUserIDFromContext(ctx, "UserID")
	logger := logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"method":  "GetTasks",
	})
	tasks, pagination, err := t.taskUsecase.GetTasks(ctx, userID, &dto.GetTasksDTO{
		Limit:   in.Limit,
		Page:    in.Page,
		Keyword: in.Keyword,
		Sort:    in.Sort,
	})
	if err != nil {
		logger.WithError(err).Error("get tasks failed")
		return nil, err
	}
	var protoTasks []*core.Task
	for _, task := range tasks {
		protoTasks = append(protoTasks, task.ToProto())
	}
	return &pbTasks.GetTasksResponse{
		Data:     protoTasks,
		Metadata: pagination.ToProto(),
	}, nil
}

func (t *TaskServiceServer) GetTaskByID(ctx context.Context, in *pbTasks.GetTaskByIDRequest) (*pbTasks.GetTaskByIDResponse, error) {
	userID, _ := util.GetUserIDFromContext(ctx, "UserID")
	logger := logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"method":  "GetTaskByID",
	})
	task, err := t.taskUsecase.GetTaskByID(ctx, userID, in.Id)
	if err != nil {
		logger.WithError(err).Error("get task by id failed")
		return nil, err
	}
	return &pbTasks.GetTaskByIDResponse{
		Data: task.ToProto(),
	}, nil
}

func (t *TaskServiceServer) CreateTask(ctx context.Context, in *pbTasks.CreateTaskRequest) (*pbTasks.CreateTaskResponse, error) {
	userID, _ := util.GetUserIDFromContext(ctx, "UserID")
	logger := logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"method":  "CreateTask",
	})
	task := &models.Task{
		UserID:      userID,
		Name:        in.Name,
		Description: in.Description,
	}
	err := t.taskUsecase.CreateTask(ctx, task)
	if err != nil {
		logger.WithError(err).Error("create task failed")
		return nil, err
	}
	return &pbTasks.CreateTaskResponse{
		Data: task.ToProto(),
	}, nil
}

func (t *TaskServiceServer) UpdateTask(ctx context.Context, in *pbTasks.UpdateTaskRequest) (*pbTypes.Empty, error) {
	userID, _ := util.GetUserIDFromContext(ctx, "UserID")
	logger := logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"method":  "UpdateTask",
	})
	task := &models.Task{
		UserID:      userID,
		Name:        in.Name,
		Description: in.Description,
	}
	task.ID = uint(in.Id)
	err := t.taskUsecase.UpdateTask(ctx, task)
	if err != nil {
		logger.WithError(err).Error("update task by id failed")
		return nil, err
	}
	return &pbTypes.Empty{}, nil
}

func (t *TaskServiceServer) DeleteTaskByID(ctx context.Context, in *pbTasks.DeleteTaskByIDRequest) (*pbTypes.Empty, error) {
	userID, _ := util.GetUserIDFromContext(ctx, "UserID")
	logger := logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"method":  "DeleteTaskByID",
	})
	err := t.taskUsecase.DeleteTaskByID(ctx, userID, in.Id)
	if err != nil {
		logger.WithError(err).Error("update task by id failed")
		return nil, err
	}
	return &pbTypes.Empty{}, nil
}
