// Package service contains business logic, defines core data types, and also responsible for interacting with data store.
package service

import (
	"context"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage"
)

// Service represents service layer for grpc.
type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// TODO: add encryption/decryption for user password and private data

func (s *Service) RegisterUser(ctx context.Context, user models.User) error {
	return s.storage.RegisterUser(ctx, user)
}

func (s *Service) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	return s.storage.GetUserByLogin(ctx, login)
}

func (s *Service) AddData(ctx context.Context, data models.Data) error {
	return s.storage.AddData(ctx, data)
}

func (s *Service) GetDataByUserID(ctx context.Context, userID string) ([]models.Data, error) {
	return s.storage.GetDataByUserID(ctx, userID)
}

func (s *Service) DeleteDataByDataID(ctx context.Context, dataID string) error {
	return s.storage.DeleteDataByDataID(ctx, dataID)
}
