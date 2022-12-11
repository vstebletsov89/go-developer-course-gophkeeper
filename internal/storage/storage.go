// Package storage defines and implements interface for Storage.
package storage

import (
	"context"
	"errors"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
)

var ErrorUnauthorized = errors.New("user is unauthorized")
var ErrorUserAlreadyExist = errors.New("user already exists")
var ErrorUserNotFound = errors.New("user not found")
var ErrorPrivateDataNotFound = errors.New("private data not found")

// Storage is the interface that must be implemented by specific storage.
type Storage interface {
	// RegisterUser registers new user in the service.
	RegisterUser(context.Context, *models.User) error
	// GetUserByLogin gets user data for authentication/authorization.
	GetUserByLogin(context.Context, string) (*models.User, error)
	// AddData adds private data to the current storage.
	AddData(context.Context, *models.Data) error
	// GetDataByUserID gets all private data for the current user.
	GetDataByUserID(context.Context, string) ([]*models.Data, error)
	// DeleteDataByDataID deletes private data for the current user.
	DeleteDataByDataID(context.Context, string) error
	// ReleaseStorage releases current storage.
	ReleaseStorage(context.Context) error
}
