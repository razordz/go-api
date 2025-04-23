package models

import "time"

type User struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"Josuel"`
	Email     string    `json:"email" example:"josuel@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2025-04-23T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-04-23T12:00:00Z"`
}
