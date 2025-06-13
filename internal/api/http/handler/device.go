package handler

import (
	"encoding/json"
	"iot-platform/internal/model"
	"iot-platform/internal/service"
	"net/http"
	"strconv"
	"time"
)

type CreateDeviceRequest struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	ApiKey string `json:"apiKey"`
}

type CreateDeviceResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
	Status  string `json:"status"`
}

type DeviceResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	ApiKey    string `json:"apiKey"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListDeviceResponse struct {
	Devices  []*DeviceResponse `json:"devices"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

func toUserResponse(device *model.Device) *DeviceResponse {
	return &DeviceResponse{
		Id:        device.Id,
		Name:      device.Name,
		Kind:      device.Kind,
		ApiKey:    device.ApiKey,
		CreatedAt: device.CreatedAt.Format(time.RFC3339),
		UpdatedAt: device.UpdatedAt.Format(time.RFC3339),
	}
}

type DeviceHandler struct {
	service service.DeviceService
}

func NewDeviceHandler(service service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		service: service,
	}
}

func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Kind == "" || req.ApiKey == "" {
		http.Error(w, "name, type and apiKey are required", http.StatusBadRequest)
		return
	}

	newDevice := &model.Device{
		Name:   req.Name,
		Kind:   req.Kind,
		ApiKey: req.ApiKey,
	}
	deviceId, err := h.service.CreateDevice(r.Context(), newDevice)
	if err != nil {
		http.Error(w, "failed to create device", http.StatusInternalServerError)
		return
	}

	response := CreateDeviceResponse{
		Message: "Device created successfully",
		Id:      deviceId,
		Status:  "200",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	devices, err := h.service.FetchDevices(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, "failed to fetch devices", http.StatusInternalServerError)
		return
	}

	var deviceResponses []*DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, toUserResponse(device))
	}

	response := ListDeviceResponse{
		Devices:  deviceResponses,
		Page:     page,
		PageSize: pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "device ID is required", http.StatusBadRequest)
		return
	}

	device, err := h.service.FindDeviceById(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to find device", http.StatusInternalServerError)
		return
	}

	response := toUserResponse(device)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "device ID is required", http.StatusBadRequest)
		return
	}

	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newDevice := &model.Device{
		Name:   req.Name,
		Kind:   req.Kind,
		ApiKey: req.ApiKey,
	}

	if err := h.service.UpdateDevice(r.Context(), id, newDevice); err != nil {
		http.Error(w, "failed to update device", http.StatusInternalServerError)
		return
	}

	response := toUserResponse(newDevice)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "device ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteDevice(r.Context(), id); err != nil {
		http.Error(w, "failed to delete device", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
