// Package storage defines and implements interface for Storage.
package storage

import (
	"context"
	"errors"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
)

// ErrorUnauthorized defines an error for unauthorized user.
var ErrorUnauthorized = errors.New("user is unauthorized")

// ErrorUserAlreadyExist defines an error for duplicate user.
var ErrorUserAlreadyExist = errors.New("user already exists")

// ErrorUserNotFound defines an error for unknown user.
var ErrorUserNotFound = errors.New("user not found")

// ErrorPrivateDataNotFound defines an error for unknown private data.
var ErrorPrivateDataNotFound = errors.New("private data not found")

// ErrorInvalidDataType defines an error for invalid private data.
var ErrorInvalidDataType = errors.New("private data has invalid type")

// Storage is the interface that must be implemented by specific storage.
type Storage interface {
	// RegisterUser registers new user in the service.
	RegisterUser(context.Context, models.User) error
	// GetUserByLogin gets user data for authentication/authorization.
	GetUserByLogin(context.Context, string) (models.User, error)
	// AddData adds private data to the current storage.
	AddData(context.Context, models.Data) error
	// GetDataByUserID gets all private data for the current user.
	GetDataByUserID(context.Context, string) ([]models.Data, error)
	// DeleteDataByDataID deletes private data for the current user.
	DeleteDataByDataID(context.Context, string) error
	// ReleaseStorage releases current storage.
	ReleaseStorage()
}
