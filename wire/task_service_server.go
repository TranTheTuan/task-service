//go:build wireinject
// +build wireinject

package wire

import (
	repoInterface "github.com/TranTheTuan/task-service/app/domain/repository"
	"github.com/TranTheTuan/task-service/app/domain/service"
	"github.com/TranTheTuan/task-service/app/domain/usecase"
	internalGrpc "github.com/TranTheTuan/task-service/app/infrastructure/grpc"
	"github.com/TranTheTuan/task-service/app/infrastructure/repository"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var (
	TaskRepoBind    = wire.Bind(new(repoInterface.TaskRepositoryInterface), new(*repository.TaskRepository))
	TaskServiceBind = wire.Bind(new(service.TaskServiceInterface), new(*service.TaskService))
	TaskUsecaseBind = wire.Bind(new(usecase.TaskUsecaseInterface), new(*usecase.TaskUsecase))

	providerSet = wire.NewSet(
		repository.NewTaskRepository,
		service.NewTaskService,
		usecase.NewTaskUsecase,
		internalGrpc.NewTaskServiceServer,
	)
)

func InitTaskServiceServer(orm *gorm.DB) *internalGrpc.TaskServiceServer {
	wire.Build(
		TaskRepoBind,
		TaskServiceBind,
		TaskUsecaseBind,
		providerSet,
	)
	return &internalGrpc.TaskServiceServer{}
}
