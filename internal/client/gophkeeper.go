package client

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
)

type SecretClient struct {
	service pb.GophkeeperClient
}

func NewSecretClient() *SecretClient {
	return &SecretClient{}
}

func (c *SecretClient) SetService(service pb.GophkeeperClient) {
	c.service = service
}

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

func (c *SecretClient) GetData(ctx context.Context) ([]models.Data, error) {
	request := &pb.GetDataRequest{}

	response, err := c.service.GetData(ctx, request)
	if err != nil {
		return nil, err
	}

	data := response.GetData()
	var convertedData []models.Data
	for _, secret := range data {
		log.Printf("secret from server storage %v", secret)
		convertedData = append(convertedData, models.Data{
			ID:         secret.GetDataId(),
			UserID:     "",
			DataType:   models.DataType(secret.GetDataType()),
			DataBinary: secret.GetDataBinary(),
		})
	}

	log.Debug().Msgf("Converted data : %v", convertedData)
	log.Debug().Msg("Client (GetData): done")
	return convertedData, err
}

func (c *SecretClient) DeleteData(ctx context.Context, dataID string) error {
	request := &pb.DeleteDataRequest{DataId: dataID}

	_, err := c.service.DeleteData(ctx, request)
	if err != nil {
		return err
	}

	log.Debug().Msg("Client (DeleteData): done")
	return nil
}
