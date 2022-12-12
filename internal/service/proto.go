package service

import (
	"github.com/google/uuid"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
)

func ConvertFromProtoDataToModel(data *proto.Data) models.Data {
	return models.Data{
		ID:         uuid.NewString(),
		UserID:     data.GetUserId(),
		DataType:   models.DataType(data.GetDataType()),
		DataBinary: data.GetDataBinary(),
	}
}

func ConvertFromModelToProtoData(data models.Data) *proto.Data {
	return &proto.Data{
		UserId:     data.UserID,
		DataType:   proto.DataType(data.DataType),
		DataBinary: data.DataBinary,
	}
}
