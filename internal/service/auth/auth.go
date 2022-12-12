// Package auth contains business logic, defines core data types, and also responsible for interacting with users.
package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	AccessToken = "uniqueAuthToken"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	Login  string `json:"login"`
	UserID string `json:"userID"`
}

func EncryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func IsUserAuthorized(user models.User, userDB models.User) (bool, error) {
	if user.Login != userDB.Login {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		return false, err
	}

	// user authorized
	return true, nil
}

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
