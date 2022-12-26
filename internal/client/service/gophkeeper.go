package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
)

// SecretClient represents a structure for core gophkeeper service.
type SecretClient struct {
	service pb.GophkeeperClient
}

// NewSecretClient returns an instance of SecretClient.
func NewSecretClient() *SecretClient {
	return &SecretClient{}
}

// SetService sets protobuf gophkeeper client.
func (c *SecretClient) SetService(service pb.GophkeeperClient) {
	c.service = service
}

// AddData is a wrapper for AddData request.
func (c *SecretClient) AddData(ctx context.Context, data models.Data) error {
	request := &pb.AddDataRequest{
		Data: &pb.Data{
			DataId:     data.ID,
			DataType:   pb.DataType(data.DataType),
			DataBinary: data.DataBinary,
		},
	}

	_, err := c.service.AddData(ctx, request)
	if err != nil {
		return err
	}

	log.Debug().Msg("Client (AddData): done")
	return nil
}

// GetData is a wrapper for GetData request.
func (c *SecretClient) GetData(ctx context.Context) ([]models.Data, error) {
	request := &pb.GetDataRequest{}

	response, err := c.service.GetData(ctx, request)
	if err != nil {
		return nil, err
	}

	data := response.GetData()
	var convertedData []models.Data
	for _, secret := range data {
		convertedData = append(convertedData, models.Data{
			ID:         secret.GetDataId(),
			UserID:     "",
			DataType:   models.DataType(secret.GetDataType()),
			DataBinary: secret.GetDataBinary(),
		})
	}

	log.Debug().Msg("Client (GetData): done")
	return convertedData, err
}

// DeleteData is a wrapper for DeleteData request.
func (c *SecretClient) DeleteData(ctx context.Context, dataID string) error {
	request := &pb.DeleteDataRequest{DataId: dataID}

	_, err := c.service.DeleteData(ctx, request)
	if err != nil {
		return err
	}

	log.Debug().Msg("Client (DeleteData): done")
	return nil
}
