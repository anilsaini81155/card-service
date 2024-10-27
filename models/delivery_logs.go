package models

import (
	"database/sql"
	"time"
)

type DeliveryLog struct {
	ID         int       `json:"id"`
	DeliveryID int       `json:"delivery_id"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateDeliveryLog creates a new log entry for a delivery.
func CreateDeliveryLog(db *sql.DB, deliveryID int, comment string) error {
	_, err := db.Exec("INSERT INTO delivery_logs (delivery_id, comment, created_at) VALUES (?, ?, ?)", deliveryID, comment, time.Now())
	return err
}
