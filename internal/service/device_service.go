package service

import (
	"context"
	"iot-platform/internal/model"
	"iot-platform/internal/repository"
	"time"
)

type deviceService interface {
	CreateDevice(ctx context.Context, device *model.Device) (string, error)
	UpdateDevice(ctx context.Context, id string, newDevice *model.Device) error
	FindDeviceById(ctx context.Context, id string) (*model.Device, error)
	FetchDevices(ctx context.Context, page int, pageSize int) ([]*model.Device, error)
	DeleteDevice(ctx context.Context, id string) error
}

type DeviceService struct {
	repo repository.DevicesRepository
}

func NewDevicesService(repo repository.DevicesRepository) *DeviceService {
	return &DeviceService{
		repo: repo,
	}
}

func (de *DeviceService) CreateDevice(ctx context.Context, device *model.Device) (string, error) {
	deviceId, err := de.repo.SaveDevice(ctx, device)
	if err != nil {
		return "", err
	}

	return deviceId, nil
}

func (de *DeviceService) UpdateDevice(ctx context.Context, id string, newDevice *model.Device) error {
	device, err := de.repo.FindDeviceById(ctx, id)
	if err != nil {
		return err
	}

	if newDevice.Name != "" {
		device.Name = newDevice.Name
	}
	if newDevice.ApiKey != "" {
		device.ApiKey = newDevice.ApiKey
	}
	device.UpdatedAt = time.Now()

	_, err = de.repo.SaveDevice(ctx, device)
	if err != nil {
		return err
	}

	return nil
}

func (de *DeviceService) FindDeviceById(ctx context.Context, id string) (*model.Device, error) {
	device, err := de.repo.FindDeviceById(ctx, id)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (de *DeviceService) FetchDevices(ctx context.Context, page int, pageSize int) ([]*model.Device, error) {
	devices, err := de.repo.ListDevices(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (de *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	err := de.repo.DeleteDevice(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
