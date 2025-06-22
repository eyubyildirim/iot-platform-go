package service

import (
	"context"
	"encoding/json"
	"iot-platform/internal/model"
	"iot-platform/internal/repository"
	"log"

	"github.com/google/uuid"
)

type ruleEvaluationService interface {
	Evaluate(ctx context.Context, data model.SensorData)
}

type RuleEvaluationService struct {
	alertRuleRepo  repository.AlertRuleRepository
	alertsTrigRepo repository.AlertsTriggeredRepository
}

func NewRuleEvaluationService(alertRuleRepo repository.AlertRuleRepository, alertsTrigRepo repository.AlertsTriggeredRepository) RuleEvaluationService {
	return RuleEvaluationService{
		alertRuleRepo:  alertRuleRepo,
		alertsTrigRepo: alertsTrigRepo,
	}
}

func (ru *RuleEvaluationService) Evaluate(ctx context.Context, data model.SensorData) {
	rules, err := ru.alertRuleRepo.ListActiveRulesForDevice(ctx, data.DeviceId)
	if err != nil {
		log.Printf("could not fetch rules for device with id: %s, err: %s", data.DeviceId, err)
		return
	}

	if len(rules) == 0 {
		log.Printf("no rules for device with id: %s", data.DeviceId)
		return
	}

	for _, rule := range rules {
		var def model.RuleDefiniton
		if err := json.Unmarshal(rule.RuleDefinition, &def); err != nil {
			log.Printf("could not process rule definition for rule with id: %s, err: %s", rule.Id, err)
			continue
		}

		if def.MetricName != data.MetricName {
			continue
		}

		isTriggered := false
		switch def.Condition {
		case "gt":
			isTriggered = data.MetricValue > def.MetricValue
		case "lt":
			isTriggered = data.MetricValue < def.MetricValue
		case "eq":
			isTriggered = data.MetricValue == def.MetricValue
		}

		if isTriggered {
			ru.createTriggeredAlert(ctx, rule, data)
		}
	}
}

func (ru *RuleEvaluationService) createTriggeredAlert(ctx context.Context, rule *model.AlertRule, data model.SensorData) {
	details, _ := json.Marshal(map[string]any{
		"ruleName":       rule.Name,
		"ruleDefinition": rule.RuleDefinition,
	})

	alert := &model.AlertTriggered{
		Id:                   uuid.NewString(),
		RuleId:               rule.Id,
		DeviceId:             data.DeviceId,
		TriggeredMetricValue: data.MetricValue,
		Details:              details,
	}

	if err := ru.alertsTrigRepo.Create(ctx, alert); err != nil {
		log.Printf("error creating the triggered alert for device with id: %s", data.DeviceId)
	}
}
