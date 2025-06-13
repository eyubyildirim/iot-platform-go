package sensordata_test

import (
	"context"
	"errors"
	"iot-platform/internal/database/postgres/sensordata"
	"iot-platform/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestSensorDataPostgresRepository_SaveSensorData_InsertSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testSensorData := &model.SensorData{
		DeviceId:    "test-device-id",
		MetricName:  "Test Name",
		MetricValue: 0.0,
	}

	mock.ExpectExec(`^INSERT INTO sensor_data \(device_id, metric_name, metric_value\) VALUES \(\$1, \$2, \$3\)$`).
		WithArgs(testSensorData.DeviceId, testSensorData.MetricName, testSensorData.MetricValue). // Arguments: ID, Name, Kind, ApiKey
		WillReturnResult(sqlmock.NewResult(0, 1))                                                 // Simulate 1 row inserted, 1 row affected (ID is not auto-increment here)

	ctx := context.Background()
	err = repo.SaveSensorData(ctx, testSensorData)

	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_SaveSensorData_InsertDbFailure(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testSensorData := &model.SensorData{
		DeviceId:    "test-device-id",
		MetricName:  "Test Name",
		MetricValue: 0.0,
	}

	mock.ExpectExec(`^INSERT INTO sensor_data \(device_id, metric_name, metric_value\) VALUES \(\$1, \$2, \$3\)$`).
		WithArgs(testSensorData.DeviceId, testSensorData.MetricName, testSensorData.MetricValue). // Arguments: ID, Name, Kind, ApiKey
		WillReturnError(errors.New("database insert error"))                                      // Simulate 1 row inserted, 1 row affected (ID is not auto-increment here)

	ctx := context.Background()
	err = repo.SaveSensorData(ctx, testSensorData)

	if err == nil {
		t.Error("expected error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_SaveSensorData_ArgumentFailure(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testSensorData := &model.SensorData{
		DeviceId:   "test-device-id",
		MetricName: "Test Name",
	}

	ctx := context.Background()
	err = repo.SaveSensorData(ctx, testSensorData)

	if err == nil {
		t.Error("expected argument error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataById_QueryDbError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := int64(1)

	mock.ExpectQuery(`^SELECT device_id, metric_name, metric_value, timestamp FROM sensor_data WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnError(errors.New("query db error"))

	ctx := context.Background()
	_, err = repo.FindSensorDataById(ctx, testId)

	if err.Error() != "query db error" {
		t.Errorf("expected query db error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataById_InvalidId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	var testId int64

	ctx := context.Background()
	_, err = repo.FindSensorDataById(ctx, testId)

	if err.Error() != "invalid id error" {
		t.Error("expected invalid id error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataById_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := int64(1)
	testRows := mock.NewRows([]string{"device_id", "metric_name", "metric_value", "timestamp"})

	mock.ExpectQuery(`^SELECT device_id, metric_name, metric_value, timestamp FROM sensor_data WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnRows(testRows)

	ctx := context.Background()
	_, err = repo.FindSensorDataById(ctx, testId)

	if err.Error() != "not found error" {
		t.Errorf("expected no error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := int64(1)
	testRows := mock.NewRows([]string{"device_id", "metric_name", "metric_value", "timestamp"})
	testRows.AddRow(uuid.NewString(), "test-metric", 1.0, time.Now())

	mock.ExpectQuery(`^SELECT device_id, metric_name, metric_value, timestamp FROM sensor_data WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnRows(testRows)

	ctx := context.Background()
	_, err = repo.FindSensorDataById(ctx, testId)

	if err != nil {
		t.Errorf("expected no error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataByDeviceId_QueryDbError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testDeviceId := "test-device-id"

	mock.ExpectQuery(`^SELECT id, metric_name, metric_value, timestamp FROM sensor_data WHERE device_id = \$1$`).
		WithArgs(testDeviceId).
		WillReturnError(errors.New("query db error"))

	ctx := context.Background()
	_, err = repo.FindSensorDataByDeviceId(ctx, testDeviceId)

	if err.Error() != "query db error" {
		t.Errorf("expected query db error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataByDeviceId_InvalidDeviceId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	var testDeviceId string

	ctx := context.Background()
	_, err = repo.FindSensorDataByDeviceId(ctx, testDeviceId)

	if err.Error() != "invalid device id error" {
		t.Error("expected invalid device id error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataByDeviceId_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testDeviceId := "test-device-id"
	testRows := mock.NewRows([]string{"id", "metric_name", "metric_value", "timestamp"})

	mock.ExpectQuery(`^SELECT id, metric_name, metric_value, timestamp FROM sensor_data WHERE device_id = \$1$`).
		WithArgs(testDeviceId).
		WillReturnRows(testRows)

	ctx := context.Background()
	_, err = repo.FindSensorDataByDeviceId(ctx, testDeviceId)

	if err.Error() != "not found error" {
		t.Errorf("expected not found error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_FindSensorDataByDeviceId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testDeviceId := "test-device-id"
	testRows := mock.NewRows([]string{"id", "metric_name", "metric_value", "timestamp"})
	testRows.AddRow(1, "test-metric", 1.0, time.Now())

	mock.ExpectQuery(`^SELECT id, metric_name, metric_value, timestamp FROM sensor_data WHERE device_id = \$1$`).
		WithArgs(testDeviceId).
		WillReturnRows(testRows)

	ctx := context.Background()
	_, err = repo.FindSensorDataByDeviceId(ctx, testDeviceId)

	if err != nil {
		t.Errorf("expected no error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_DeleteSensorData_QueryDbError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := int64(1)

	mock.ExpectExec(`^DELETE FROM sensor_data WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnError(errors.New("query db error"))

	ctx := context.Background()
	err = repo.DeleteSensorData(ctx, testId)

	if err.Error() != "query db error" {
		t.Errorf("expected query db error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_DeleteSensorData_InvalidId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	var testId int64

	ctx := context.Background()
	err = repo.DeleteSensorData(ctx, testId)

	if err.Error() != "invalid id error" {
		t.Error("expected invalid id error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_DeleteSensorData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := int64(1)

	mock.ExpectExec(`^DELETE FROM sensor_data WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate 1 row deleted

	ctx := context.Background()
	err = repo.DeleteSensorData(ctx, testId)

	if err != nil {
		t.Errorf("expected no error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_ListSensorData_QueryDbError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testPage := 1
	testPageSize := 10

	mock.ExpectQuery(`^SELECT id, device_id, metric_name, metric_value, timestamp FROM sensor_data LIMIT \$1 OFFSET \$2$`).
		WithArgs(testPageSize, (testPage-1)*testPageSize).
		WillReturnError(errors.New("query db error"))

	ctx := context.Background()
	_, err = repo.ListSensorData(ctx, testPage, testPageSize)

	if err.Error() != "query db error" {
		t.Errorf("expected query db error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_ListSensorData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testPage := 1
	testPageSize := 10
	testRows := mock.NewRows([]string{"id", "device_id", "metric_name", "metric_value", "timestamp"})
	testRows.AddRow(1, "test-device-id", "test-metric", 1.0, time.Now())

	mock.ExpectQuery(`^SELECT id, device_id, metric_name, metric_value, timestamp FROM sensor_data LIMIT \$1 OFFSET \$2$`).
		WithArgs(testPageSize, (testPage-1)*testPageSize).
		WillReturnRows(testRows)

	ctx := context.Background()
	_, err = repo.ListSensorData(ctx, testPage, testPageSize)

	if err != nil {
		t.Errorf("expected no error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_ListSensorData_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testPage := 1
	testPageSize := 10
	testRows := mock.NewRows([]string{"id", "device_id", "metric_name", "metric_value", "timestamp"})

	mock.ExpectQuery(`^SELECT id, device_id, metric_name, metric_value, timestamp FROM sensor_data LIMIT \$1 OFFSET \$2$`).
		WithArgs(testPageSize, (testPage-1)*testPageSize).
		WillReturnRows(testRows)

	ctx := context.Background()
	_, err = repo.ListSensorData(ctx, testPage, testPageSize)

	if err.Error() != "not found error" {
		t.Errorf("expected not found error, but got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_ListSensorData_InvalidPage(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	var testPage int
	testPageSize := 10

	ctx := context.Background()
	_, err = repo.ListSensorData(ctx, testPage, testPageSize)

	if err.Error() != "invalid page error" {
		t.Error("expected invalid page error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSensorDataPostgresRepository_ListSensorData_InvalidPageSize(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testPage := 1
	var testPageSize int

	ctx := context.Background()
	_, err = repo.ListSensorData(ctx, testPage, testPageSize)

	if err.Error() != "invalid page size error" {
		t.Error("expected invalid page size error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
