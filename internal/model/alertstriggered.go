package model

import (
	"encoding/json"
	"time"
)

type AlertTriggered struct {
	Id                   string          `json:"id"`
	RuleId               string          `json:"ruleId"`
	DeviceId             string          `json:"deviceId"`
	TriggeredMetricValue float64         `json:"triggeredMetricValue"`
	Details              json.RawMessage `json:"details"`
	TriggeredAt          time.Time       `json:"triggeredAt"`
}
