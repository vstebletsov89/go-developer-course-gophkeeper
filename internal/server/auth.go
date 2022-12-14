package server

import (
	"context"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServer
	service service.Service
	jwt     auth.JWT
}

func NewAuthServer(service service.Service, jwt auth.JWT) *AuthServer {
	return &AuthServer{service: service, jwt: jwt}
}

func (a *AuthServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var response pb.RegisterResponse

	encryptedPassword, err := auth.EncryptPassword(request.GetUser().GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := models.User{
		ID:       uuid.NewString(),
		Login:    request.GetUser().GetLogin(),
		Password: encryptedPassword,
	}

	err = a.service.RegisterUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("Server (Register): done")
	return &response, nil
}

func (a *AuthServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse

	userDB, err := a.service.GetUserByLogin(ctx, request.GetUser().GetLogin())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}
	user := models.User{
		ID:       "",
		Login:    request.GetUser().GetLogin(),
		Password: request.GetUser().GetPassword(),
	}

	ok, err := auth.IsUserAuthorized(&user, &userDB)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect username/password")
	}

	token, err := a.jwt.GenerateToken(userDB.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	response.User = request.GetUser()
	response.Token = &pb.Token{
		UserId: userDB.ID,
		Token:  token,
	}

	log.Debug().Msg("Server (Login): done")
	return &response, nil
}

// AuthFunc is used by a middleware to authenticate requests.
func (a *AuthServer) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := a.jwt.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", tokenInfo)

	newCtx := context.WithValue(ctx, auth.AccessToken, tokenInfo)

	return newCtx, nil
}
