package notifier

import "iot-platform/internal/model"

type Notifier interface {
	Send(alert model.AlertTriggered, rule model.AlertRule) error
}

type NotificationPayload struct {
	Alert model.AlertTriggered
	Rule  model.AlertRule
}
