package usecase

import (
	"context"
	"fmt"

	"github.com/TranTheTuan/task-service/app/domain/dto"
	"github.com/TranTheTuan/task-service/app/infrastructure/grpc/client"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUseCase struct {
	authClient client.AuthClient
}

func (a *AuthUseCase) AuthHandler(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	authID, err := a.authClient.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	serverTransportStream := grpc.ServerTransportStreamFromContext(ctx)
	pass, err := a.authClient.Authorize(ctx, &dto.AuthorizeDTO{
		CasbinUser: fmt.Sprint(authID),
		RequestURI: "",
		Method: serverTransportStream.Method(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Enforce error: %v", err)
	}
	if !pass {
		return nil, nil
	}

	ctx = context.WithValue(ctx, "UserIDKey", authID)

	return ctx, nil
}
