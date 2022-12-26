// Package auth contains business logic, defines core data types, and also responsible for interacting with users.
package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// UserContextType user context type.
type UserContextType string

const (
	// UserCtx defines user context name.
	UserCtx UserContextType = "UserCtx"
)

// JWTManager represents a structure for jwt manager.
type JWTManager struct {
	secretKey string
}

// UserClaims custom claims for jwt.
type UserClaims struct {
	jwt.RegisteredClaims
}

// JWT interface is the interface that must be implemented by JWTManager.
type JWT interface {
	GenerateToken(user string) (string, error)
	ValidateToken(token string) (*UserClaims, error)
}

// NewJWTManager return an instance of JWTManager.
func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey: secretKey}
}

// check that JWTManager implements all required methods.
var _ JWT = (*JWTManager)(nil)

// GenerateToken generates jwt token.
func (j *JWTManager) GenerateToken(user string) (string, error) {
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
	genToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	log.Debug().Msgf("GenerateToken token: %s", genToken)
	return genToken, nil
}

// ValidateToken verifies that jwt token is valid.
func (j *JWTManager) ValidateToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(j.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	log.Debug().Msg("ValidateToken: OK")
	return claims, nil
}

// EncryptPassword is used for encryption user password.
func EncryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// IsUserAuthorized is used to compare current user with registered user in storage.
func IsUserAuthorized(user *models.User, userDB *models.User) (bool, error) {
	if user.Login != userDB.Login {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		return false, err
	}

	log.Debug().Msg("User authorized")
	// user authorized
	return true, nil
}

// ExtractUserIDFromContext extracts userID from context.
func ExtractUserIDFromContext(ctx context.Context) string {
	// try to get userID from context
	userID, ok := ctx.Value(UserCtx).(string)
	if ok {
		log.Debug().Msgf("ExtractUserIDFromContext: '%s'", userID)
		return userID
	}
	return ""
}
