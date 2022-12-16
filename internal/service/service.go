// Package service contains business logic, defines core data types, and also responsible for interacting with data store.
package service

import (
	"context"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage"
)

// Service represents service layer for grpc server.
type Service struct {
	storage storage.Storage
}

// NewService returns an instance of Service.
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// RegisterUser is a wrapper for storage layer. It is used in grpc server methods.
func (s *Service) RegisterUser(ctx context.Context, user models.User) error {
	return s.storage.RegisterUser(ctx, user)
}

// GetUserByLogin is a wrapper for storage layer. It is used in grpc server methods.
func (s *Service) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	return s.storage.GetUserByLogin(ctx, login)
}

// AddData is a wrapper for storage layer. It is used in grpc server methods.
func (s *Service) AddData(ctx context.Context, data models.Data) error {
	return s.storage.AddData(ctx, data)
}

// GetDataByUserID is a wrapper for storage layer. It is used in grpc server methods.
func (s *Service) GetDataByUserID(ctx context.Context, userID string) ([]models.Data, error) {
	return s.storage.GetDataByUserID(ctx, userID)
}

// DeleteDataByDataID is a wrapper for storage layer. It is used in grpc server methods.
func (s *Service) DeleteDataByDataID(ctx context.Context, dataID string) error {
	return s.storage.DeleteDataByDataID(ctx, dataID)
}
