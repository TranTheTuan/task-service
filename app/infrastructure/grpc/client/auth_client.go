package client

import (
	"context"

	pbAuth "github.com/TranTheTuan/pbtypes/build/go/auth"
	"github.com/TranTheTuan/task-service/app/domain/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	authorizeClient pbAuth.AuthAuthorizeServiceClient
	verifyClient pbAuth.AuthVerifyServiceClient
}

func NewAuthClient(authClientAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(authClientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// defer conn.Close()

	authorizeClient := pbAuth.NewAuthAuthorizeServiceClient(conn)
	verifyClient := pbAuth.NewAuthVerifyServiceClient(conn)
	return &AuthClient{authorizeClient, verifyClient}, nil
}

func (a *AuthClient) Authorize(ctx context.Context, in *dto.AuthorizeDTO) (bool, error) {
	res, err := a.authorizeClient.Authorize(ctx, &pbAuth.AuthorizeRequest{
		CasbinUser: in.CasbinUser,
		RequestUri: in.RequestURI,
		Method:     in.Method,
	})
	if err != nil {
		return false, err
	}
	return res.Pass, nil
}

func (a *AuthClient) VerifyToken(ctx context.Context, token string) (uint32, error) {
	res, err := a.verifyClient.VerifyToken(ctx, &pbAuth.VerifyTokenRequest{Token: token})
	if err != nil {
		return 0, err
	}
	return res.Id, nil
}
