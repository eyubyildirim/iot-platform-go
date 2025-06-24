package alertstriggered

import (
	"context"
	"database/sql"
	"fmt"
	"iot-platform/internal/model"
	"log"
)

type AlertsTriggeredPostgresRepository struct {
	db *sql.DB
}

func NewAlertsTriggeredPostgresRepository(db *sql.DB) *AlertsTriggeredPostgresRepository {
	return &AlertsTriggeredPostgresRepository{db: db}
}

func (r *AlertsTriggeredPostgresRepository) Create(ctx context.Context, alert *model.AlertTriggered) error {
	query := `INSERT INTO alerts_triggered (id, rule_id, device_id, triggered_value, details) VALUES ($1, $2, $3, $4, $5)`

	log.Printf("new alert trigger: %s", alert.Details)

	_, err := r.db.ExecContext(ctx, query,
		alert.Id,
		alert.RuleId,
		alert.DeviceId,
		alert.TriggeredMetricValue,
		alert.Details,
	)

	if err != nil {
		return fmt.Errorf("failed to insert triggered alert: %w", err)
	}

	return nil
}
