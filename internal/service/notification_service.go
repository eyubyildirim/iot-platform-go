package service

import (
	"iot-platform/internal/api/notifier"
	"iot-platform/internal/model"
	"log"
)

type NotificationService struct {
	DefaultNotifier notifier.Notifier
}

func NewNotificationService(notifier notifier.Notifier) *NotificationService {
	return &NotificationService{
		DefaultNotifier: notifier,
	}
}

func (ns *NotificationService) Dispatch(alert model.AlertTriggered, rule model.AlertRule) {
	if err := ns.DefaultNotifier.Send(alert, rule); err != nil {
		log.Printf("error sending notification: %s", err.Error())
	}
}
