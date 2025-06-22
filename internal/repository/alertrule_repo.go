package repository

import (
	"context"
	"iot-platform/internal/model"
)

type AlertRuleRepository interface {
	CreateRule(ctx context.Context, alertRule *model.AlertRule) (string, error)
	FindRuleByID(ctx context.Context, id string) (*model.AlertRule, error)
	DeleteRule(ctx context.Context, id string) error
	ListByDeviceId(ctx context.Context, deviceId string, page, pageSize int) ([]*model.AlertRule, error)
	ListActiveRules(ctx context.Context) ([]*model.AlertRule, error)
	UpdateRule(ctx context.Context, alertRule *model.AlertRule) error
	ToggleRuleStatus(ctx context.Context, id string) (bool, error)
}
