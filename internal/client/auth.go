package client

import (
	"context"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
)

type AuthClient struct {
	user        models.User
	accessToken string
	service     pb.AuthClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{}
}

func (a *AuthClient) SetService(service pb.AuthClient) {
	a.service = service
}

func (a *AuthClient) AccessToken() string {
	return a.accessToken
}

func (a *AuthClient) SetAccessToken(accessToken string) {
	a.accessToken = accessToken
}

func (a *AuthClient) User() models.User {
	return a.user
}

func (a *AuthClient) SetUser(user models.User) {
	a.user = user
}

func (a *AuthClient) Login(ctx context.Context) (string, error) {
	request := &pb.LoginRequest{
		User: &pb.User{
			Login:    a.user.Login,
			Password: a.user.Password,
		},
	}

	response, err := a.service.Login(ctx, request)
	if err != nil {
		return "", err
	}

	a.user.ID = response.GetToken().GetUserId()
	log.Debug().Msgf("Client (Login): done %v", a.user)
	return response.GetToken().GetToken(), nil
}

func (a *AuthClient) Register(ctx context.Context) error {
	request := &pb.RegisterRequest{
		User: &pb.User{
			Login:    a.user.Login,
			Password: a.user.Password,
		},
	}

	_, err := a.service.Register(ctx, request)
	if err != nil {
		return err
	}

	log.Debug().Msg("Client (Register): done")
	return nil
}

func (a *AuthClient) UnaryInterceptorClient(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(map[string]string{auth.AccessToken: a.AccessToken(),
		auth.UserCtx: a.User().ID})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
