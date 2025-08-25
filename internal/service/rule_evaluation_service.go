package service

import (
	"context"
	"encoding/json"
	"iot-platform/internal/api/notifier"
	"iot-platform/internal/model"
	"iot-platform/internal/repository"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RuleState struct {
	BreachStartTime time.Time
	AlertSent       bool
}

type stateKey struct {
	RuleId   string
	DeviceId string
}

type ruleEvaluationService interface {
	Evaluate(ctx context.Context, data model.SensorData)
}

type RuleEvaluationService struct {
	alertRuleRepo  repository.AlertRuleRepository
	alertsTrigRepo repository.AlertsTriggeredRepository
	stateCache     map[stateKey]RuleState
	mu             sync.Mutex
	notifChannel   chan<- notifier.NotificationPayload
}

func NewRuleEvaluationService(alertRuleRepo repository.AlertRuleRepository, alertsTrigRepo repository.AlertsTriggeredRepository, notifChannel chan<- notifier.NotificationPayload) RuleEvaluationService {
	return RuleEvaluationService{
		alertRuleRepo:  alertRuleRepo,
		alertsTrigRepo: alertsTrigRepo,
		stateCache:     make(map[stateKey]RuleState),
		notifChannel:   notifChannel,
	}
}

func (ru *RuleEvaluationService) Evaluate(ctx context.Context, data model.SensorData) {
	log.Printf("evaluating rules for metric name: %s", data.MetricName)
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
		key := stateKey{
			RuleId:   rule.Id,
			DeviceId: data.DeviceId,
		}

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
		case "between":
			isTriggered = data.MetricValue >= def.MetricValueMin && data.MetricValue <= def.MetricValueMax
		}

		ru.mu.Lock()

		_, isTracked := ru.stateCache[key]

		if isTriggered {
			if !isTracked {
				log.Printf("condition met for rule %s, starting to track", rule.Name)
				ru.stateCache[key] = RuleState{
					BreachStartTime: time.Now(),
					AlertSent:       false,
				}
			}

			state := ru.stateCache[key]
			durationReq := time.Duration(def.DurationMinutes) * time.Minute

			if !state.AlertSent && time.Since(state.BreachStartTime) >= durationReq {
				log.Printf("condition has been met for rule '%s' by device with id: %s, alert is sent\n", rule.Name, data.DeviceId)
				ru.createTriggeredAlert(ctx, rule, data)

				state.AlertSent = true
				ru.stateCache[key] = state
			}
		} else {
			if isTracked {
				log.Printf("condition for rule '%s' with device id '%s' has cleared, resetting.\n", rule.Name, data.DeviceId)
				delete(ru.stateCache, key)
			}
		}

		ru.mu.Unlock()
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
		TriggeredAt:          time.Now(),
		Details:              details,
	}

	if err := ru.alertsTrigRepo.Create(ctx, alert); err != nil {
		log.Printf("error creating the triggered alert for device with id: %s", data.DeviceId)
		return
	}

	ru.notifChannel <- notifier.NotificationPayload{
		Alert: *alert,
		Rule:  *rule,
	}
}
