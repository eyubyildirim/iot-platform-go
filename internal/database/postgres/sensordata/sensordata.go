package sensordata

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iot-platform/internal/model"
	"log"
)

func InitDb(host, port, user, pass, dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

type SensorDataPostgresRepository struct {
	db *sql.DB
}

func (se *SensorDataPostgresRepository) SaveSensorData(ctx context.Context, sensorData *model.SensorData) error {
	if sensorData.DeviceId == "" || sensorData.MetricName == "" {
		return errors.New("save argument error")
	}

	_, err := se.db.Exec("INSERT INTO sensor_data (device_id, metric_name, metric_value) VALUES ($1, $2, $3)", sensorData.DeviceId, sensorData.MetricName, sensorData.MetricValue)
	if err != nil {
		return err
	}

	return nil
}

func (se *SensorDataPostgresRepository) FindSensorDataById(ctx context.Context, id int64) (*model.SensorData, error) {
	if id == 0 {
		return nil, errors.New("invalid id error")
	}

	rows, err := se.db.Query("SELECT device_id, metric_name, metric_value, timestamp FROM sensor_data WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensorData model.SensorData
	if rows.Next() {
		err = rows.Scan(&sensorData.DeviceId, &sensorData.MetricName, &sensorData.MetricValue, &sensorData.Timestamp)
		if err != nil {
			return nil, errors.New("scan error")
		}
	}

	if sensorData.DeviceId == "" {
		return nil, errors.New("not found error")
	}

	return &sensorData, nil
}

func (se *SensorDataPostgresRepository) FindSensorDataByDeviceId(ctx context.Context, deviceId string) ([]*model.SensorData, error) {
	if deviceId == "" {
		return nil, errors.New("invalid device id error")
	}

	rows, err := se.db.Query("SELECT id, metric_name, metric_value, timestamp FROM sensor_data WHERE device_id = $1", deviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensorDataList []*model.SensorData
	for rows.Next() {
		var sensorData model.SensorData
		err = rows.Scan(&sensorData.Id, &sensorData.MetricName, &sensorData.MetricValue, &sensorData.Timestamp)
		if err != nil {
			return nil, errors.New("scan error")
		}
		sensorDataList = append(sensorDataList, &sensorData)
	}

	if len(sensorDataList) == 0 {
		return nil, errors.New("not found error")
	}

	return sensorDataList, nil
}

func (se *SensorDataPostgresRepository) DeleteSensorData(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("invalid id error")
	}

	result, err := se.db.Exec("DELETE FROM sensor_data WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("not found error")
	}

	return nil
}

func (se *SensorDataPostgresRepository) ListSensorData(ctx context.Context, page int, pageSize int) ([]*model.SensorData, error) {
	if page <= 0 {
		return nil, errors.New("invalid page error")
	}
	if pageSize <= 0 {
		return nil, errors.New("invalid page size error")
	}

	offset := (page - 1) * pageSize
	rows, err := se.db.Query("SELECT id, device_id, metric_name, metric_value, timestamp FROM sensor_data LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensorDataList []*model.SensorData
	for rows.Next() {
		var sensorData model.SensorData
		err = rows.Scan(&sensorData.Id, &sensorData.DeviceId, &sensorData.MetricName, &sensorData.MetricValue, &sensorData.Timestamp)
		if err != nil {
			return nil, errors.New("scan error")
		}
		sensorDataList = append(sensorDataList, &sensorData)
	}

	if len(sensorDataList) == 0 {
		return nil, errors.New("not found error")
	}

	return sensorDataList, nil
}

func NewSensorDataPostgresRepository(db *sql.DB) (*SensorDataPostgresRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}

	return &SensorDataPostgresRepository{
		db: db,
	}, nil
}
