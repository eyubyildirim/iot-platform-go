package repository

import (
	"context"
	"iot-platform/internal/model"
)

// AlertsTriggeredRepository defines the contract for storing triggered alert records.
type AlertsTriggeredRepository interface {
	Create(ctx context.Context, alert *model.AlertTriggered) error
}
