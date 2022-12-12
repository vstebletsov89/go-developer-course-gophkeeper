package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

const (
	AccessToken = "uniqueAuthToken"
)

func ExtractUserIDFromContext(ctx context.Context) string {
	// try to get userID from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get(AccessToken)
		if len(values) > 0 {
			userID := values[0]
			log.Debug().Msgf("ExtractUserIDFromContext (GRPC): '%s'", userID)
			return userID
		}
	}
	return "" // token not found
}
