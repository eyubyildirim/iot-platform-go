package repository

import (
	"context"
	"iot-platform/internal/model"
)

type DevicesRepository interface {
	SaveDevice(ctx context.Context, device *model.Device) error
	FindDeviceById(ctx context.Context, id string) (*model.Device, error)
	DeleteDevice(ctx context.Context, id string) error
	ListDevices(ctx context.Context, page, pageSize int) ([]*model.Device, error)
}
