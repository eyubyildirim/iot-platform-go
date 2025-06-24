package main

import (
	"context"
	"iot-platform/internal/api/http/handler"
	"iot-platform/internal/api/http/middleware"
	"iot-platform/internal/api/notifier"
	"iot-platform/internal/api/notifier/notifiers"
	"iot-platform/internal/database/postgres"
	"iot-platform/internal/database/postgres/alertrule"
	"iot-platform/internal/database/postgres/alertstriggered"
	"iot-platform/internal/database/postgres/device"
	"iot-platform/internal/database/postgres/sensordata"
	"iot-platform/internal/model"
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

	dataProcessingChannel := make(chan model.SensorData, 1000)
	notificationChannel := make(chan notifier.NotificationPayload, 100)

	deviceRepo, _ := device.NewDevicePostgresRepository(db)
	sensorDataRepo, _ := sensordata.NewSensorDataPostgresRepository(db)
	alertRuleRepo, _ := alertrule.NewAlertRulesPostgresRepository(db)
	alertsTrigRepo := alertstriggered.NewAlertsTriggeredPostgresRepository(db)

	deviceService := service.NewDevicesService(deviceRepo)
	sensorDataService := service.NewSensorDataService(sensorDataRepo, dataProcessingChannel)
	alertRuleService := service.NewAlertRuleService(alertRuleRepo)
	ruleEvalService := service.NewRuleEvaluationService(alertRuleRepo, alertsTrigRepo, notificationChannel)
	notifService := service.NewNotificationService(notifiers.LogNotifier{})

	log.Println("Starting rule evaluation worker...")
	go func() {
		// This loop will run forever, waiting for data to arrive on the channel.
		for data := range dataProcessingChannel {
			// Each piece of data is evaluated in turn.
			ruleEvalService.Evaluate(context.Background(), data)
		}
	}()

	log.Println("Starting notification dispatcher worker...")
	go func() {
		for payload := range notificationChannel {
			notifService.Dispatch(payload.Alert, payload.Rule)
		}
	}()

	deviceHandler := handler.NewDeviceHandler(*deviceService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /devices", deviceHandler.ListDevices)
	mux.HandleFunc("POST /devices", deviceHandler.CreateDevice)
	mux.HandleFunc("GET /devices/{id}", deviceHandler.GetDevice)
	mux.HandleFunc("PUT /devices/{id}", deviceHandler.UpdateDevice)
	mux.HandleFunc("DELETE /devices/{id}", deviceHandler.DeleteDevice)

	sensorDataHandler := handler.NewSensorDataHandler(*sensorDataService)
	mux.HandleFunc("GET /sensor-data", sensorDataHandler.ListSensorData)
	mux.Handle("POST /sensor-data/{deviceId}", middleware.AuthMiddleware(deviceRepo)(http.HandlerFunc(sensorDataHandler.CreateSensorData)))
	mux.HandleFunc("GET /sensor-data/{id}", sensorDataHandler.GetSensorDataByDeviceId)
	mux.HandleFunc("DELETE /sensor-data/{id}", sensorDataHandler.DeleteSensorData)

	alertRuleHandler := handler.NewAlertRuleHandler(*alertRuleService)
	mux.HandleFunc("POST /alert-rule", alertRuleHandler.CreateRule)
	mux.HandleFunc("GET /alert-rule/{id}", alertRuleHandler.GetRuleByID)
	mux.HandleFunc("PATCH /alert-rule/toggle-status/{id}", alertRuleHandler.ToggleRuleStatus)
	mux.HandleFunc("DELETE /alert-rule/{id}", alertRuleHandler.DeleteRule)
	mux.HandleFunc("GET /alert-rule/active", alertRuleHandler.ListActiveRules)
	mux.HandleFunc("GET /alert-rule/device/{id}", alertRuleHandler.ListByDeviceId)
	mux.HandleFunc("PATCH /alert-rule/{id}", alertRuleHandler.UpdateRule)

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
