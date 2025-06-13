package main

import (
	"iot-platform/internal/api/http/handler"
	"iot-platform/internal/database/postgres"
	"iot-platform/internal/database/postgres/device"
	"iot-platform/internal/database/postgres/sensordata"
	"iot-platform/internal/service"
	"log"
	"net/http"
	"time"
)

func main() {
	config, err := loadConfiguration("/Users/eyubyildirim/Documents/go-projects/iot-platform/config.json")
	if err != nil {
		log.Fatalf("problem parsing config: %s", err)
	}

	db, err := postgres.InitDb(config.Database.Host, config.Database.Port, config.Database.User, config.Database.Pass, config.Database.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deviceRepo, err := device.NewDevicePostgresRepository(db)
	if err != nil {
		log.Fatal("error connecting to database")
	}
	deviceService := service.NewDevicesService(deviceRepo)

	sensorDataRepo, err := sensordata.NewSensorDataPostgresRepository(db)
	if err != nil {
		log.Fatal("error connecting to database")
	}
	sensorDataService := service.NewSensorDataService(sensorDataRepo)

	deviceHandler := handler.NewDeviceHandler(*deviceService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /devices", deviceHandler.ListDevices)
	mux.HandleFunc("POST /devices", deviceHandler.CreateDevice)
	mux.HandleFunc("GET /devices/{id}", deviceHandler.GetDevice)
	mux.HandleFunc("PUT /devices/{id}", deviceHandler.UpdateDevice)
	mux.HandleFunc("DELETE /devices/{id}", deviceHandler.DeleteDevice)

	sensorDataHandler := handler.NewSensorDataHandler(*sensorDataService)
	mux.HandleFunc("GET /sensor-data", sensorDataHandler.ListSensorData)
	mux.HandleFunc("POST /sensor-data", sensorDataHandler.CreateSensorData)
	mux.HandleFunc("GET /sensor-data/{id}", sensorDataHandler.GetSensorDataByDeviceId)
	mux.HandleFunc("DELETE /sensor-data/{id}", sensorDataHandler.DeleteSensorData)

	server := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Server starting on port %s\n", config.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("Server stopped gracefully")
}
