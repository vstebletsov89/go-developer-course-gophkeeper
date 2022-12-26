// Package service contains service layer for grpc clients.
package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthClient represents a structure for authorization service.
type AuthClient struct {
	user        models.User
	accessToken string
	service     pb.AuthClient
}

// NewAuthClient returns an instance of AuthClient.
func NewAuthClient() *AuthClient {
	return &AuthClient{}
}

// SetService sets protobuf authorization client.
func (a *AuthClient) SetService(service pb.AuthClient) {
	a.service = service
}

// AccessToken getter for accessToken.
func (a *AuthClient) AccessToken() string {
	return a.accessToken
}

// SetAccessToken setter for accessToken.
func (a *AuthClient) SetAccessToken(accessToken string) {
	a.accessToken = accessToken
}

// User getter for current user.
func (a *AuthClient) User() models.User {
	return a.user
}

// SetUser setter for current user.
func (a *AuthClient) SetUser(user models.User) {
	a.user = user
}

// Login is a wrapper for Login request.
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

// Register is a wrapper for Register request.
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

// UnaryInterceptorClient is a client interceptor for attaching access token and current userID.
func (a *AuthClient) UnaryInterceptorClient(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+a.AccessToken())
	log.Debug().Msgf("UnaryInterceptorClient (attaching bearer with jwt token): %v", a.AccessToken())
	return invoker(newCtx, method, req, reply, cc, opts...)
}
