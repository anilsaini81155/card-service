package models

import (
	"database/sql"
	"errors"
	"time"
)

type Delivery struct {
	ID        int       `json:"id"`
	CardID    string    `json:"card_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateDelivery creates a new delivery entry.
func CreateDelivery(db *sql.DB, cardID int, status string) (int64, error) {
	result, err := db.Exec("INSERT INTO deliveries (card_id, status, created_at, updated_at) VALUES (?, ?, ?, ?)", cardID, status, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateDeliveryStatus updates the status of a delivery entry.
func UpdateDeliveryStatus(db *sql.DB, deliveryID int, status string) error {
	_, err := db.Exec("UPDATE deliveries SET status = ?, updated_at = ? WHERE id = ?", status, time.Now(), deliveryID)
	return err
}

// GetDeliveryByCardID fetches the delivery status by card ID.
func GetDeliveryByCardID(db *sql.DB, cardID int) (*Delivery, error) {
	delivery := &Delivery{}
	err := db.QueryRow("SELECT id, card_id, status FROM deliveries WHERE card_id = ?", cardID).Scan(
		&delivery.ID, &delivery.CardID, &delivery.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("delivery not found")
		}
		return nil, err
	}

	return delivery, nil
}
