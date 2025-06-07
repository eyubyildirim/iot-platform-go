package model

import "time"

type Device struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Kind      string    `json:"type"`
	ApiKey    string    `json:"apiKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
