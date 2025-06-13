package service

import (
	"context"
	"iot-platform/internal/model"
	"iot-platform/internal/repository"
)

type sensorDataService interface {
	CreateSensorData(ctx context.Context, sensorData *model.SensorData) error
	FindSensorDataById(ctx context.Context, id int64) (*model.SensorData, error)
	FindSensorDataByDeviceId(ctx context.Context, deviceId string) ([]*model.SensorData, error)
	FetchSensorData(ctx context.Context, page int, pageSize int) ([]*model.SensorData, error)
	DeleteSensorData(ctx context.Context, id int64) error
}

type SensorDataService struct {
	repo repository.SensorDataRepository
}

func NewSensorDataService(repo repository.SensorDataRepository) *SensorDataService {
	return &SensorDataService{
		repo: repo,
	}
}

func (se *SensorDataService) CreateSensorData(ctx context.Context, sensorData *model.SensorData) error {
	err := se.repo.SaveSensorData(ctx, sensorData)
	if err != nil {
		return err
	}

	return nil
}

func (se *SensorDataService) FindSensorDataById(ctx context.Context, id int64) (*model.SensorData, error) {
	sensorData, err := se.repo.FindSensorDataById(ctx, id)
	if err != nil {
		return nil, err
	}

	return sensorData, nil
}

func (se *SensorDataService) FindSensorDataByDeviceId(ctx context.Context, deviceId string) ([]*model.SensorData, error) {
	sensorData, err := se.repo.FindSensorDataByDeviceId(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	return sensorData, nil
}

func (se *SensorDataService) FetchSensorData(ctx context.Context, page int, pageSize int) ([]*model.SensorData, error) {
	sensorData, err := se.repo.ListSensorData(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return sensorData, nil
}

func (se *SensorDataService) DeleteSensorData(ctx context.Context, id int64) error {
	err := se.repo.DeleteSensorData(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

