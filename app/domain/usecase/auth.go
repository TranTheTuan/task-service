package usecase

import (
	"context"
	"fmt"

	"github.com/TranTheTuan/task-service/app/domain/dto"
	"github.com/TranTheTuan/task-service/app/infrastructure/grpc/client"
	"github.com/sirupsen/logrus"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthUsecase struct {
	authClient *client.AuthClient
}

func NewAuthUsecase(authClient *client.AuthClient) *AuthUsecase {
	return &AuthUsecase{
		authClient,
	}
}

func (a *AuthUsecase) AuthHandler(ctx context.Context) (context.Context, error) {
	logger := logrus.WithFields(logrus.Fields{
		"method": "AuthHandler",
	})
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		logger.WithError(err).Error("get token failed")
		return nil, err
	}

	userID, err := a.authClient.VerifyToken(ctx, token)
	if err != nil {
		logger.WithError(err).Error("verify token failed")
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	pattern := md.Get("pattern")
	if len(pattern) == 0 {
		pattern = []string{""}
	}
	method := md.Get("method")
	if len(method) == 0 {
		method = []string{""}
	}
	logger.WithFields(logrus.Fields{
		"userID":  userID,
		"pattern": pattern[0],
		"method":  method[0],
	}).Info("get pattern successfully")
	_, err = a.authClient.Authorize(ctx, &dto.AuthorizeDTO{
		CasbinUser: fmt.Sprint(userID),
		RequestURI: pattern[0],
		Method:     method[0],
	})
	if err != nil {
		logger.WithError(err).Error("authorize failed")
		return nil, status.Errorf(codes.PermissionDenied, "Authorize error: %v", err)
	}
	// if !pass {
	// 	logger.Error("unauthorized")
	// 	return nil, errors.New("unauthorized")
	// }

	ctx = context.WithValue(ctx, "UserID", userID)

	return ctx, nil
}
