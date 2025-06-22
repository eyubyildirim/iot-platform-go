package handler

import (
	"encoding/json"
	"iot-platform/internal/database/postgres/alertrule"
	"iot-platform/internal/model"
	"iot-platform/internal/service"
	"log"
	"net/http"
	"time"
)

type CreateAlertRuleRequest struct {
	Name           string              `json:"name"`
	DeviceId       *string             `json:"deviceId,omitempty"`
	RuleDefinition model.RuleDefiniton `json:"ruleDefinition"`
	IsActive       bool                `json:"isActive"`
}

type UpdateAlertRuleRequest struct {
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	DeviceId       *string             `json:"deviceId,omitempty"`
	RuleDefinition model.RuleDefiniton `json:"ruleDefinition"`
	IsActive       bool                `json:"isActive"`
}

type AlertRuleResponse struct {
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	DeviceId       *string             `json:"deviceId,omitempty"`
	RuleDefinition model.RuleDefiniton `json:"ruleDefinition"`
	IsActive       bool                `json:"isActive"`
	CreatedAt      string              `json:"createdAt"`
	UpdatedAt      string              `json:"updatedAt"`
}

type AlertRuleListResponse struct {
	TotalCount int                 `json:"totalCount"`
	AlertRules []AlertRuleResponse `json:"alertRules"`
}

func toAlertRuleResponse(alertRule *model.AlertRule) (*AlertRuleResponse, error) {
	var ruleDef model.RuleDefiniton
	err := json.Unmarshal(alertRule.RuleDefinition, &ruleDef)
	if err != nil {
		return nil, err
	}

	return &AlertRuleResponse{
		Id:             alertRule.Id,
		Name:           alertRule.Name,
		DeviceId:       alertRule.DeviceId,
		RuleDefinition: ruleDef,
		IsActive:       alertRule.IsActive,
		CreatedAt:      alertRule.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      alertRule.UpdatedAt.Format(time.RFC3339),
	}, nil
}

type AlertRuleHandler struct {
	service service.AlertRuleService
}

func NewAlertRuleHandler(service service.AlertRuleService) *AlertRuleHandler {
	return &AlertRuleHandler{
		service: service,
	}
}

func (h *AlertRuleHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req CreateAlertRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ruleDefBytes, err := json.Marshal(req.RuleDefinition)

	alertRule := &model.AlertRule{
		Name:           req.Name,
		DeviceId:       req.DeviceId,
		RuleDefinition: ruleDefBytes,
		IsActive:       req.IsActive,
	}

	id, err := h.service.CreateRule(r.Context(), alertRule)
	if err != nil {
		http.Error(w, "Failed to create alert rule", http.StatusInternalServerError)
		return
	}

	alertRule.Id = id

	response, err := toAlertRuleResponse(alertRule)
	if err != nil {
		http.Error(w, "Failed to convert alert rule to response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AlertRuleHandler) GetRuleByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing alert rule ID", http.StatusBadRequest)
		return
	}

	alertRule, err := h.service.FindRuleByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve alert rule", http.StatusInternalServerError)
		return
	}
	if alertRule == nil {
		http.Error(w, "Alert rule not found", http.StatusNotFound)
		return
	}

	response, err := toAlertRuleResponse(alertRule)
	if err != nil {
		http.Error(w, "Failed to convert alert rule to response", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AlertRuleHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing alert rule ID", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteRule(r.Context(), id)
	if err != nil {
		if err == alertrule.ErrRuleNotFound {
			http.Error(w, "Alert rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete alert rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AlertRuleHandler) ListByDeviceId(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("id")
	if deviceId == "" {
		http.Error(w, "Missing device ID", http.StatusBadRequest)
		return
	}

	page := 1
	pageSize := 10

	alertRules, err := h.service.ListByDeviceId(r.Context(), deviceId, page, pageSize)
	if err != nil {
		http.Error(w, "Failed to list alert rules", http.StatusInternalServerError)
		return
	}

	var response AlertRuleListResponse
	response.TotalCount = len(alertRules)
	for _, rule := range alertRules {
		responseRule, err := toAlertRuleResponse(rule)
		if err != nil {
			// Log the error and skip this rule
			log.Printf("Warning: could not process rule %s: %v", rule.Id, err)
			continue
		}
		response.AlertRules = append(response.AlertRules, *responseRule)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AlertRuleHandler) ListActiveRules(w http.ResponseWriter, r *http.Request) {
	alertRules, err := h.service.ListActiveRules(r.Context())
	if err != nil {
		http.Error(w, "Failed to list active alert rules", http.StatusInternalServerError)
		return
	}

	var response AlertRuleListResponse
	response.TotalCount = len(alertRules)
	for _, rule := range alertRules {
		responseRule, err := toAlertRuleResponse(rule)
		if err != nil {
			// Log the error and skip this rule
			log.Printf("Warning: could not process rule %s: %v", rule.Id, err)
			continue
		}
		response.AlertRules = append(response.AlertRules, *responseRule)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AlertRuleHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	var req UpdateAlertRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ruleDefBytes, err := json.Marshal(req.RuleDefinition)
	if err != nil {
		http.Error(w, "Invalid rule definition", http.StatusBadRequest)
		return
	}

	alertRule := &model.AlertRule{
		Id:             req.Id,
		Name:           req.Name,
		DeviceId:       req.DeviceId,
		RuleDefinition: ruleDefBytes,
		IsActive:       req.IsActive,
	}

	err = h.service.UpdateRule(r.Context(), alertRule)
	if err != nil {
		http.Error(w, "Failed to update alert rule", http.StatusInternalServerError)
		return
	}

	response, err := toAlertRuleResponse(alertRule)
	if err != nil {
		http.Error(w, "Failed to convert alert rule to response", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func (h *AlertRuleHandler) ToggleRuleStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing alert rule ID", http.StatusBadRequest)
		return
	}

	isActive, err := h.service.ToggleRuleStatus(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to toggle alert rule status", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{"isActive": isActive}
	json.NewEncoder(w).Encode(response)
}
