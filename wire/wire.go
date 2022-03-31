//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/TranTheTuan/task-service/app/domain/usecase"
	"github.com/TranTheTuan/task-service/app/infrastructure/grpc/client"
)

func InitAuthUsecase(authGrpcAddr string) (*usecase.AuthUsecase, error) {
	wire.Build(
		client.NewAuthClient,
		usecase.NewAuthUsecase,
	)

	return &usecase.AuthUsecase{}, nil
}
