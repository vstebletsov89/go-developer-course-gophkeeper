package server

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type JwtInterceptor struct {
	jwt auth.JWT
}

func NewJwtInterceptor(jwt auth.JWT) *JwtInterceptor {
	return &JwtInterceptor{jwt: jwt}
}

func (j *JwtInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if strings.Contains(info.FullMethod, "Register") || strings.Contains(info.FullMethod, "Login") {
		// skip validation jwt token for register and login
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is empty")
	}

	values := md.Get(auth.AccessToken)
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "jwt token not provided")
	}

	token := values[0]
	err := j.jwt.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token: %v", err)
	}

	log.Debug().Msg("Interceptor authorization: OK")
	return handler(ctx, req)
}
