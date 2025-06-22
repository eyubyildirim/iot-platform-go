package model

import (
	"encoding/json"
	"time"
)

type RuleDefiniton struct {
	MetricName      string  `json:"metricName"`
	Condition       string  `json:"condition"`
	MetricValue     float64 `json:"metricValue"`
	DurationMinutes int     `json:"durationMinutes"`
}

type AlertRule struct {
	Id             string          `json:"id"`
	Name           string          `json:"name"`
	DeviceId       *string         `json:"deviceId,omitempty"`
	RuleDefinition json.RawMessage `json:"ruleDefinition"`
	IsActive       bool            `json:"isActive"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}
