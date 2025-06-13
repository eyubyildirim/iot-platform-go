package main

import (
	"context"
	"iot-platform/internal/database/postgres/device"
	"iot-platform/internal/model"
	"iot-platform/internal/service"
	"log"
)

func main() {
	config, err := loadConfiguration("/Users/eyubyildirim/Documents/go-projects/iot-platform/config.json")
	if err != nil {
		log.Fatalf("problem parsing config: %s", err)
	}

	db, err := device.InitDb(config.Database.Host, config.Database.Port, config.Database.User, config.Database.Pass, config.Database.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo, err := device.NewDevicePostgresRepository(db)
	if err != nil {
		log.Fatal("error connecting to database")
	}
	service := service.NewDevicesService(repo)

	err = service.CreateDevice(context.Background(), &model.Device{
		Name:   "TestDevice2",
		Kind:   "TestType",
		ApiKey: "1111aaaa",
	})
	if err != nil {
		log.Fatal(err)
	}
}
