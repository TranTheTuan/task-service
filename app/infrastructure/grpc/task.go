package grpc

import "context"

import (
	pbTasks "github.com/TranTheTuan/pbtypes/build/go/tasks"
	pbTypes "google.golang.org/protobuf/types/known/emptypb"
)

type TaskServiceServer struct {

}

func NewTaskServiceServer() *TaskServiceServer {
	return &TaskServiceServer{}
}

func (t *TaskServiceServer) GetTasks(ctx context.Context, _ *pbTypes.Empty) (*pbTasks.GetTasksResponse, error) {
	return nil, nil
}

func (t *TaskServiceServer) GetTaskByID(ctx context.Context, in *pbTasks.GetTaskByIDRequest) (*pbTasks.GetTaskByIDResponse, error) {
	return nil, nil
}

func (t *TaskServiceServer) CreateTask(ctx context.Context, in *pbTasks.CreateTaskRequest) (*pbTasks.CreateTaskResponse, error) {
	return nil, nil
}

func (t *TaskServiceServer) UpdateTask(ctx context.Context, in *pbTasks.UpdateTaskRequest) (*pbTypes.Empty, error) {
	return nil, nil
}

func (t *TaskServiceServer) DeleteTaskByID(ctx context.Context, in *pbTasks.DeleteTaskByIDRequest) (*pbTypes.Empty, error) {
	return nil, nil
}
