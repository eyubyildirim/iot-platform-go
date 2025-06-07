package database_test

import (
	"context"
	"errors"
	"fmt"
	"iot-platform/internal/database"
	"iot-platform/internal/model"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock" // Import sqlmock
	"github.com/google/uuid"
)

func TestDevicePostgresRepository_SaveDevice_InsertSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testDevice := &model.Device{
		Name:   "Error Device",
		Kind:   "Error",
		ApiKey: "error-api-key",
		Id:     "",
	}

	mock.ExpectExec(`^INSERT INTO devices \(id, name, kind, api_key\) VALUES \(\$1, \$2, \$3, \$4\)$`).
		WithArgs(sqlmock.AnyArg(), testDevice.Name, testDevice.Kind, testDevice.ApiKey). // Arguments: ID, Name, Kind, ApiKey
		WillReturnResult(sqlmock.NewResult(1, 1))                                        // Simulate 1 row inserted, 1 row affected (ID is not auto-increment here)

	ctx := context.Background()
	err = repo.SaveDevice(ctx, testDevice)

	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_SaveDevice_InsertFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testDevice := &model.Device{
		Name:   "Failed Insert Name",
		Kind:   "Failed Insert Kind",
		ApiKey: "failed-insert-api-key",
		Id:     "",
	}
	insertErr := errors.New("failed to insert")

	mock.ExpectExec(`^INSERT INTO devices \(id, name, kind, api_key\) VALUES \(\$1, \$2, \$3, \$4\)$`).
		WithArgs(sqlmock.AnyArg(), testDevice.Name, testDevice.Kind, testDevice.ApiKey).
		WillReturnError(insertErr)

	err = repo.SaveDevice(context.Background(), testDevice)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_SaveDevice_InsertArgumentError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testDevice := &model.Device{
		Name:   "",
		Kind:   "Failed Insert Kind",
		ApiKey: "failed-insert-api-key",
		Id:     "",
	}

	err = repo.SaveDevice(context.Background(), testDevice)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_SaveDevice_UpdateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.New()
	testDevice := &model.Device{
		Name:   "Success name",
		Kind:   "Success",
		ApiKey: "success-api-key",
		Id:     id.String(),
	}

	expectedSql := `^UPDATE devices SET name = \$1, kind = \$2, api_key = \$3, updated_at = \$4 WHERE id = \$5$`

	mock.ExpectExec(expectedSql).
		WithArgs(testDevice.Name, testDevice.Kind, testDevice.ApiKey, sqlmock.AnyArg(), testDevice.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err = repo.SaveDevice(ctx, testDevice)

	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_SaveDevice_UpdateFailure(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.NewString()
	testDevice := &model.Device{
		Name:   "Failure Update Name",
		Kind:   "Failure Update Kind",
		ApiKey: "failure-update-api-key",
		Id:     id,
	}

	expectedSql := `^UPDATE devices SET name = \$1, kind = \$2, api_key = \$3, updated_at = \$4 WHERE id = \$5$`
	updateErr := errors.New("failed to update")

	mock.ExpectExec(expectedSql).
		WithArgs(testDevice.Name, testDevice.Kind, testDevice.ApiKey, sqlmock.AnyArg(), testDevice.Id).
		WillReturnError(updateErr)

	ctx := context.Background()
	err = repo.SaveDevice(ctx, testDevice)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_NewDevicePostgresRepository_PingFails(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pingErr := errors.New("failed db connection")
	mock.ExpectPing().WillReturnError(pingErr)

	_, err = database.NewDevicePostgresRepository(db)

	if err == nil {
		log.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_FindDeviceById_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := "Not Found Device Id"
	notFoundErr := errors.New("device not found error")

	mock.ExpectQuery(`^SELECT devices.id, devices.name, devices.kind, devices.api_key, devices.created_at, devices.updated_at FROM devices WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnError(notFoundErr)

	_, err = repo.FindDeviceById(context.Background(), testId)

	if err == nil {
		t.Error("expected not found error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_FindDeviceById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := uuid.NewString()

	rows := sqlmock.NewRows([]string{"id", "name", "kind", "api_key", "created_at", "updated_at"})
	rows.AddRow(testId, "Success Name", "Success Kind", "success-api-key", time.Now(), time.Now())

	mock.ExpectQuery(`^SELECT devices.id, devices.name, devices.kind, devices.api_key, devices.created_at, devices.updated_at FROM devices WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnRows(rows)

	testDevice, err := repo.FindDeviceById(context.Background(), testId)

	fmt.Println(testDevice.Id)
	if err != nil {
		t.Error("expected not found error, got nil")
	}

	if testDevice.Id != testId {
		t.Error("expected test device is not the right one")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_DeleteDevice_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := "Not Found Device Id"

	mock.ExpectExec(`^DELETE FROM devices WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteDevice(context.Background(), testId)

	if err == nil {
		t.Error("expected not found error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_DeleteDevice_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := "Success Id"

	mock.ExpectExec(`^DELETE FROM devices WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteDevice(context.Background(), testId)

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_DeleteDevice_DeleteDbError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := "Success Id"
	dbErr := errors.New("db error")

	mock.ExpectExec(`^DELETE FROM devices WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnError(dbErr)

	err = repo.DeleteDevice(context.Background(), testId)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_DeleteDevice_ResultError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testId := "Success Id"
	resErr := errors.New("result error")

	mock.ExpectExec(`^DELETE FROM devices WHERE id = \$1$`).
		WithArgs(testId).
		WillReturnResult(sqlmock.NewErrorResult(resErr))

	err = repo.DeleteDevice(context.Background(), testId)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_ListDevices_EmptyList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	testPage := 1
	testPageSize := 10
	testRows := sqlmock.NewRows([]string{"id", "name", "kind", "api_key", "created_at", "updated_at"})

	mock.ExpectQuery(`^SELECT devices.id, devices.name, devices.kind, devices.api_key, devices.created_at, devices.updated_at FROM devices ORDER BY created_at OFFSET \$1 LIMIT \$2$`).
		WithArgs((testPage-1)*testPageSize, testPageSize).
		WillReturnRows(testRows)

	_, err = repo.ListDevices(context.Background(), testPage, testPageSize)

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_ListDevices_SelectDbError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	dbErr := errors.New("db error")

	mock.ExpectQuery(`^SELECT devices.id, devices.name, devices.kind, devices.api_key, devices.created_at, devices.updated_at FROM devices ORDER BY created_at OFFSET \$1 LIMIT \$2$`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(dbErr)

	_, err = repo.ListDevices(context.Background(), 1, 10)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDevicePostgresRepository_ListDevices_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	dbErr := errors.New("db error")
	testRows := sqlmock.NewRows([]string{"id", "name", "kind", "api_key", "created_at", "updated_at"})
	testRows.AddRow("Read Error Id", "Read Error Name", "Read Error Kind", "read-error-api-key", time.Now(), time.Now())
	testRows.RowError(0, dbErr)

	mock.ExpectQuery(`^SELECT devices.id, devices.name, devices.kind, devices.api_key, devices.created_at, devices.updated_at FROM devices ORDER BY created_at OFFSET \$1 LIMIT \$2$`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(testRows)

	_, err = repo.ListDevices(context.Background(), 1, 10)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
