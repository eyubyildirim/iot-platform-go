package notifiers

import (
	"iot-platform/internal/model"
	"log"
	"time"
)

type LogNotifier struct{}

func (lo LogNotifier) Send(alert model.AlertTriggered, rule model.AlertRule) error {
	log.Println("--- 🚀 NOTIFICATION DISPATCHED 🚀 ---")
	log.Printf("Alert For: '%s'\n", rule.Name)
	log.Printf("Device ID: %s\n", alert.DeviceId)
	log.Printf("Triggered Value: %.2f\n", alert.TriggeredMetricValue)
	log.Printf("Timestamp: %s\n", alert.TriggeredAt.Format(time.RFC1123))
	log.Println("------------------------------------")
	return nil
}
