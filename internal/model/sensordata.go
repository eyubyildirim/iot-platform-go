package model

import "time"

type SensorData struct {
	Id          int64     `json:"id"`
	DeviceId    string    `json:"deviceId"`
	MetricName  string    `json:"metricName"`
	MetricValue float64   `json:"metricValue"`
	Timestamp   time.Time `json:"timestamp"`
}
