package repository

import (
	"context"
	"iot-platform/internal/model"
)

type SensorDataRepository interface {
	SaveSensorData(ctx context.Context, sensorData *model.SensorData) error
	FindSensorDataById(ctx context.Context, id int64) (*model.SensorData, error)
	FindSensorDataByDeviceId(ctx context.Context, id string) ([]*model.SensorData, error)
	DeleteSensorData(ctx context.Context, id int64) error
	ListSensorData(ctx context.Context, page, pageSize int) ([]*model.SensorData, error)
}
