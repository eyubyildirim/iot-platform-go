package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"iot-platform/internal/model"
	"iot-platform/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SensorDataHandler struct {
	sensorDataService service.SensorDataService
}

type CreateSensorDataRequest struct {
	DeviceId    string  `json:"deviceId"`
	MetricName  string  `json:"metricName"`
	Metricvalue float64 `json:"metricValue"`
}

type CreateSensorDataResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DeleteSensorDataResponse CreateSensorDataResponse

type SensorDataResponse struct {
	Id          int64   `json:"id"`
	DeviceId    string  `json:"deviceId"`
	MetricName  string  `json:"metricName"`
	MetricValue float64 `json:"metricValue"`
	Timestamp   string  `json:"timestamp"`
}

type ListSensorDataResponse struct {
	SensorData []*SensorDataResponse `json:"sensorData"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
}

func toSensorData(device *model.SensorData) *SensorDataResponse {
	return &SensorDataResponse{
		Id:          device.Id,
		DeviceId:    device.DeviceId,
		MetricName:  device.MetricName,
		MetricValue: device.MetricValue,
		Timestamp:   device.Timestamp.Format(time.RFC3339),
	}
}
func NewSensorDataHandler(sensorDataService service.SensorDataService) *SensorDataHandler {
	return &SensorDataHandler{
		sensorDataService: sensorDataService,
	}
}

func (h *SensorDataHandler) CreateSensorData(w http.ResponseWriter, r *http.Request) {
	var request CreateSensorDataRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if request.DeviceId == "" || request.MetricName == "" || request.Metricvalue <= 0 {
		http.Error(w, "Device ID, Metric Name and Metric Value are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	sensorData := &model.SensorData{
		DeviceId:    request.DeviceId,
		MetricName:  request.MetricName,
		MetricValue: request.Metricvalue,
		Timestamp:   time.Now(),
	}

	if err := h.sensorDataService.CreateSensorData(ctx, sensorData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create sensor data: %v", err), http.StatusInternalServerError)
		return
	}

	response := CreateSensorDataResponse{
		Message: "Sensor data created successfully",
		Status:  "success",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SensorDataHandler) ListSensorData(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	sensorDataList, err := h.sensorDataService.FetchSensorData(ctx, pageInt, pageSizeInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch sensor data: %v", err), http.StatusInternalServerError)
		return
	}

	response := ListSensorDataResponse{
		SensorData: make([]*SensorDataResponse, len(sensorDataList)),
		Page:       pageInt,
		PageSize:   pageSizeInt,
	}
	for i, sensorData := range sensorDataList {
		response.SensorData[i] = toSensorData(sensorData)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Fetched %d sensor data records for page %d with page size %d", len(sensorDataList), pageInt, pageSizeInt)
}

func (h *SensorDataHandler) GetSensorDataByDeviceId(w http.ResponseWriter, r *http.Request) {
	deviceId := strings.Split(r.URL.Path, "/")[2] // Assuming the URL is like /sensor-data/{deviceId}

	if deviceId == "" {
		http.Error(w, "Device ID is required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	sensorDataList, err := h.sensorDataService.FindSensorDataByDeviceId(ctx, deviceId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get sensor data: %v", err), http.StatusInternalServerError)
		return
	}

	response := ListSensorDataResponse{
		SensorData: make([]*SensorDataResponse, len(sensorDataList)),
		Page:       1,
		PageSize:   len(sensorDataList),
	}
	for i, sensorData := range sensorDataList {
		response.SensorData[i] = toSensorData(sensorData)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Fetched %d sensor data records for device ID %s", len(sensorDataList), deviceId)
}

func (h *SensorDataHandler) DeleteSensorData(w http.ResponseWriter, r *http.Request) {
	sensorDataId := strings.Split(r.URL.Path, "/")[2] // Assuming the URL is like /sensor-data/{id}
	if sensorDataId == "" {
		http.Error(w, "Sensor Data ID is required", http.StatusBadRequest)
		return
	}

	sensorDataIdInt, err := strconv.ParseInt(sensorDataId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Sensor Data ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.sensorDataService.DeleteSensorData(ctx, sensorDataIdInt); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete sensor data: %v", err), http.StatusInternalServerError)
		return
	}
	response := DeleteSensorDataResponse{
		Message: "Sensor data deleted successfully",
		Status:  "success",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Sensor data with ID %s deleted successfully", sensorDataId)
}
