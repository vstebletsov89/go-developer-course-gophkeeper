package server

import (
	"context"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"google.golang.org/grpc"
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

	// TODO: implemt token validation and add it to bearer
	//userID := uuid.NewString()
	//
	//if md, ok := metadata.FromIncomingContext(ctx); ok {
	//	values := md.Get(service.AccessToken)
	//	if len(values) > 0 {
	//		userID = values[0]
	//		log.Printf("UnaryInterceptor userID from context: '%s'", userID)
	//	}
	//}
	//
	//md, ok := metadata.FromIncomingContext(ctx)
	//if ok {
	//	md.Append(service.AccessToken, string(userID))
	//}
	//newCtx := metadata.NewIncomingContext(ctx, md)

	return handler(ctx, req)
	//return handler(newCtx, req)
}
