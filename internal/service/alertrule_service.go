package service

import (
	"context"
	"iot-platform/internal/model"
	"iot-platform/internal/repository"
)

type alertRuleService interface {
	CreateRule(ctx context.Context, alertRule *model.AlertRule) (string, error)
	FindRuleByID(ctx context.Context, id string) (*model.AlertRule, error)
	DeleteRule(ctx context.Context, id string) error
	ListByDeviceId(ctx context.Context, deviceId string, page, pageSize int) ([]*model.AlertRule, error)
	ListActiveRules(ctx context.Context) ([]*model.AlertRule, error)
	UpdateRule(ctx context.Context, alertRule *model.AlertRule) error
	ToggleRuleStatus(ctx context.Context, id string) (bool, error)
}

type AlertRuleService struct {
	repo repository.AlertRuleRepository
}

func NewAlertRuleService(repo repository.AlertRuleRepository) *AlertRuleService {
	return &AlertRuleService{
		repo: repo,
	}
}

func (al *AlertRuleService) CreateRule(ctx context.Context, alertRule *model.AlertRule) (string, error) {
	id, err := al.repo.CreateRule(ctx, alertRule)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (al *AlertRuleService) FindRuleByID(ctx context.Context, id string) (*model.AlertRule, error) {
	alertRule, err := al.repo.FindRuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if alertRule == nil {
		return nil, nil // No rule found
	}
	return alertRule, nil
}

func (al *AlertRuleService) DeleteRule(ctx context.Context, id string) error {
	err := al.repo.DeleteRule(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (al *AlertRuleService) ListByDeviceId(ctx context.Context, deviceId string, page int, pageSize int) ([]*model.AlertRule, error) {
	alertRules, err := al.repo.ListByDeviceId(ctx, deviceId, page, pageSize)
	if err != nil {
		return nil, err
	}
	return alertRules, nil
}

func (al *AlertRuleService) ListActiveRules(ctx context.Context) ([]*model.AlertRule, error) {
	alertRules, err := al.repo.ListActiveRules(ctx)
	if err != nil {
		return nil, err
	}
	return alertRules, nil
}

func (al *AlertRuleService) UpdateRule(ctx context.Context, alertRule *model.AlertRule) error {
	err := al.repo.UpdateRule(ctx, alertRule)
	if err != nil {
		return err
	}
	return nil
}

func (al *AlertRuleService) ToggleRuleStatus(ctx context.Context, id string) (bool, error) {
	return al.repo.ToggleRuleStatus(ctx, id)
}
