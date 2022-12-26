package server

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// JwtInterceptor represents a structure for jwt interceptor.
type JwtInterceptor struct {
	jwt auth.JWT
}

// NewJwtInterceptor returns an instance of JwtInterceptor.
func NewJwtInterceptor(jwt auth.JWT) *JwtInterceptor {
	return &JwtInterceptor{jwt: jwt}
}

// UnaryInterceptor grpc interceptor to validate access token. It is used for authorization of users.
func (j *JwtInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Debug().Msg("Interceptor authorization (grpc_middleware)")

	if strings.Contains(info.FullMethod, "Register") || strings.Contains(info.FullMethod, "Login") {
		// skip validation jwt token for register and login
		return handler(ctx, req)
	}

	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Validation token: %v", token)
	claims, err := j.jwt.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token: %v", err)
	}

	newCtx := context.WithValue(ctx, auth.UserCtx, claims.ID)

	log.Debug().Msg("Interceptor authorization: OK")
	return handler(newCtx, req)
}
