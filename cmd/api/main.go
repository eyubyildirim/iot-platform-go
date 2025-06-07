package main

import (
	"context"
	"iot-platform/internal/database"
	"iot-platform/internal/model"
	"iot-platform/internal/service"
	"log"
)

func main() {
	db, err := database.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo, err := database.NewDevicePostgresRepository(db)
	if err != nil {
		log.Fatal("error connecting to database")
	}
	service := service.NewDevicesService(repo)

	err = service.CreateDevice(context.Background(), &model.Device{
		Id:     "53bafafa-31e2-4fcb-bfa6-628b6671ee19",
		Name:   "TestDevice2",
		Kind:   "TestType",
		ApiKey: "1111aaaa",
	})
	if err != nil {
		log.Fatal(err)
	}
}
