package device

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iot-platform/internal/model"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DevicePostgresRepository struct {
	db *sql.DB
}

func NewDevicePostgresRepository(db *sql.DB) (*DevicePostgresRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}

	return &DevicePostgresRepository{
		db: db,
	}, nil
}

func (de *DevicePostgresRepository) SaveDevice(ctx context.Context, device *model.Device) (string, error) {
	if device.Id == "" {
		fmt.Println("Test")
		newDeviceId := uuid.New().String()
		_, err := de.db.Exec(`INSERT INTO devices (id, name, kind, api_key) VALUES ($1, $2, $3, $4)`, newDeviceId, device.Name, device.Kind, device.ApiKey)

		if err != nil {
			return newDeviceId, err
		}

		return newDeviceId, nil
	} else {
		_, err := de.db.Exec(`UPDATE devices SET name = $1, kind = $2, api_key = $3, updated_at = $4 WHERE id = $5`, device.Name, device.Kind, device.ApiKey, time.Now(), device.Id)
		if err != nil {
			return device.Id, err
		}

		return device.Id, nil
	}

}

func (de *DevicePostgresRepository) FindDeviceById(ctx context.Context, id string) (*model.Device, error) {
	row := de.db.QueryRow(`SELECT devices.id, devices.name, devices.kind, devices.api_key, devices.created_at, devices.updated_at FROM devices WHERE id = $1`, id)

	var device model.Device

	err := row.Scan(&device.Id, &device.Name, &device.Kind, &device.ApiKey, &device.CreatedAt, &device.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (de *DevicePostgresRepository) DeleteDevice(ctx context.Context, id string) error {
	res, err := de.db.Exec(`DELETE FROM devices WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no device found with id: %s", id)
	}

	return nil
}

func (de *DevicePostgresRepository) ListDevices(ctx context.Context, page int, pageSize int) ([]*model.Device, error) {
	rows, err := de.db.Query(`SELECT devices.id, devices.name, devices.kind, devices.api_key, devices.created_at, devices.updated_at FROM devices ORDER BY created_at OFFSET $1 LIMIT $2`, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*model.Device
	for rows.Next() {
		var device model.Device
		err := rows.Scan(&device.Id, &device.Name, &device.Kind, &device.ApiKey, &device.UpdatedAt, &device.CreatedAt)
		if err != nil {
			return nil, err
		}

		devices = append(devices, &device)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return devices, nil
}
