// Package auth contains business logic, defines core data types, and also responsible for interacting with users.
package auth

import (
	"context"
	"fmt"
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
	secretKey string
}

type UserClaims struct {
	jwt.RegisteredClaims
}

type JWT interface {
	GenerateToken(user string) (string, error)
	ValidateToken(token string) (*UserClaims, error)
}

// check that JWTManager implements all required methods.
var _ JWT = (*JWTManager)(nil)

func (J *JWTManager) GenerateToken(user string) (string, error) {
	claims := UserClaims{RegisteredClaims: jwt.RegisteredClaims{
		Issuer:    "Gophkeeper",
		Subject:   "authorization",
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        user,
	}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	genToken, err := token.SignedString(J.secretKey)
	if err != nil {
		return "", err
	}

	return genToken, nil
}

func (J *JWTManager) ValidateToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return J.secretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func EncryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func IsUserAuthorized(user *models.User, userDB *models.User) (bool, error) {
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
