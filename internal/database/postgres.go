package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iot-platform/internal/model"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func InitDb() (*sql.DB, error) {
	connStr := "postgres://eyub:1234@localhost:5432/iot_platform?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

type DevicePostgresRepository struct {
	db *sql.DB
}

func NewDevicePostgresRepository(db *sql.DB) (*DevicePostgresRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed db connection")
	}

	return &DevicePostgresRepository{
		db: db,
	}, nil
}

func (de *DevicePostgresRepository) SaveDevice(ctx context.Context, device *model.Device) error {
	if device.Id == "" {
		if device.Name == "" || device.Kind == "" || device.ApiKey == "" {
			return errors.New("argument error")
		}

		query := `INSERT INTO devices (id, name, kind, api_key) VALUES ($1, $2, $3, $4)`
		fmt.Println(query)

		_, err := de.db.Exec(query, uuid.New().String(), device.Name, device.Kind, device.ApiKey)
		if err != nil {
			return errors.New("failed to insert")
		}

		return nil
	} else {
		query := "UPDATE devices SET name = $1, kind = $2, api_key = $3, updated_at = $4 WHERE id = $5"
		updatedAt := time.Now()
		fmt.Printf("App Query: '%s'\n", query)
		_, err := de.db.ExecContext(ctx, query, device.Name, device.Kind, device.ApiKey, updatedAt, device.Id)
		if err != nil {
			return errors.New("failed to update")
		}

		return nil
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

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return errors.New("device not found error")
	}

	return nil
}

func (de *DevicePostgresRepository) ListDevices(ctx context.Context, page, pageSize int) ([]*model.Device, error) {
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
