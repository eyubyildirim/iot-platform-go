package alertrule

import (
	"context"
	"database/sql"
	"errors"
	"iot-platform/internal/model"
)

type AlertRulesPostgresRepository struct {
	db *sql.DB
}

var (
	ErrRuleNotFound = errors.New("alert rule not found")
)

func NewAlertRulesPostgresRepository(db *sql.DB) (*AlertRulesPostgresRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}

	return &AlertRulesPostgresRepository{
		db: db,
	}, nil
}

func (al *AlertRulesPostgresRepository) CreateRule(ctx context.Context, alertRule *model.AlertRule) (string, error) {
	query := `INSERT INTO alert_rules (name, device_id, rule_definition, is_active) VALUES ($1, $2, $3, $4) RETURNING id, craeted_at, updated_at`

	var id string
	err := al.db.QueryRowContext(ctx, query,
		alertRule.Name,
		alertRule.DeviceId,
		alertRule.RuleDefinition,
		alertRule.IsActive).Scan(&id, &alertRule.CreatedAt, &alertRule.UpdatedAt)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (al *AlertRulesPostgresRepository) FindRuleByID(ctx context.Context, id string) (*model.AlertRule, error) {
	query := `SELECT id, name, device_id, rule_definition, is_active, created_at FROM alert_rules WHERE id = $1`
	row := al.db.QueryRowContext(ctx, query, id)

	alertRule := &model.AlertRule{}
	err := row.Scan(&alertRule.Id, &alertRule.Name, &alertRule.DeviceId, &alertRule.RuleDefinition, &alertRule.IsActive, &alertRule.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRuleNotFound // No rule found
		}
		return nil, err // Other error
	}

	return alertRule, nil
}

func (al *AlertRulesPostgresRepository) DeleteRule(ctx context.Context, id string) error {
	query := `DELETE FROM alert_rules WHERE id = $1`
	result, err := al.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRuleNotFound // No rule found to delete
	}

	return nil
}

func (al *AlertRulesPostgresRepository) ListByDeviceId(ctx context.Context, deviceId string, page int, pageSize int) ([]*model.AlertRule, error) {
	query := `SELECT id, name, device_id, rule_definition, is_active, created_at FROM alert_rules WHERE device_id = $1 LIMIT $2 OFFSET $3`
	rows, err := al.db.QueryContext(ctx, query, deviceId, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alertRules []*model.AlertRule
	for rows.Next() {
		alertRule := &model.AlertRule{}
		if err := rows.Scan(&alertRule.Id, &alertRule.Name, &alertRule.DeviceId, &alertRule.RuleDefinition, &alertRule.IsActive, &alertRule.CreatedAt); err != nil {
			return nil, err
		}
		alertRules = append(alertRules, alertRule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alertRules, nil
}

func (al *AlertRulesPostgresRepository) ListActiveRules(ctx context.Context) ([]*model.AlertRule, error) {
	query := `SELECT id, name, device_id, rule_definition, is_active, created_at FROM alert_rules WHERE is_active = true`
	rows, err := al.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alertRules []*model.AlertRule
	for rows.Next() {
		alertRule := &model.AlertRule{}
		if err := rows.Scan(&alertRule.Id, &alertRule.Name, &alertRule.DeviceId, &alertRule.RuleDefinition, &alertRule.IsActive, &alertRule.CreatedAt); err != nil {
			return nil, err
		}
		alertRules = append(alertRules, alertRule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alertRules, nil
}

func (al *AlertRulesPostgresRepository) UpdateRule(ctx context.Context, alertRule *model.AlertRule) error {
	query := `UPDATE alert_rules SET name = $1, device_id = $2, rule_definition = $3 WHERE id = $4`
	result, err := al.db.ExecContext(ctx, query,
		alertRule.Name,
		alertRule.DeviceId,
		alertRule.RuleDefinition,
		alertRule.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRuleNotFound // No rule found to update
	}

	return nil
}

func (al *AlertRulesPostgresRepository) ToggleRuleStatus(ctx context.Context, id string) (bool, error) {
	query := `UPDATE alert_rules SET is_active = NOT is_active WHERE id = $1 RETURNING is_active`

	var isActive bool
	err := al.db.QueryRowContext(ctx, query, id).Scan(&isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ErrRuleNotFound // No rule found to toggle
		}
		return false, err // Other error
	}

	return isActive, nil
}
