package models

import (
	"time"
)

type Message struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"` // "chat_message", "location_update"...
	OrderID    int64     `json:"order_id"`
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	Latitude   float64   `json:"latitude,omitempty"`  // Cho location
	Longitude  float64   `json:"longitude,omitempty"` // Cho location
	CreatedAt  time.Time `json:"created_at"`
}
