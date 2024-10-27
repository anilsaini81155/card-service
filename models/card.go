package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

type Card struct {
	ID        int       `json:"id"`
	CardNo    string    `json:"card_no"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GenerateCardNo() string {
	rand.Seed(uint64(time.Now().UnixNano())) // Explicit conversion to int64

	// rand.Seed(time.Now().UnixNano())
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	cardNo := make([]rune, 3)
	for i := range cardNo {
		cardNo[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%s%04d", string(cardNo), rand.Intn(10000))
}

// CreateCard creates a new card for a user.
func CreateCard(db *sql.DB, cardNo string, userID int) (int64, error) {
	// Ensure cardNo is unique
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM cards WHERE card_no = ?)", cardNo).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("card number already exists")
	}

	// Insert the card into the database
	result, err := db.Exec("INSERT INTO cards (card_no, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)", cardNo, userID, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetCardByID fetches a card by its ID.
func GetCardByID(db *sql.DB, cardID int) (*Card, error) {
	card := &Card{}
	err := db.QueryRow("SELECT id, card_no, user_id FROM cards WHERE id = ?", cardID).Scan(
		&card.ID, &card.CardNo, &card.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("card not found")
		}
		return nil, err
	}

	return card, nil
}

// GetCardByID fetches a card by its ID.
func GetCardByUserID(db *sql.DB, userId int) (*Card, error) {
	card := &Card{}
	err := db.QueryRow("SELECT id, card_no, user_id FROM cards WHERE user_id = ?", userId).Scan(
		&card.ID, &card.CardNo, &card.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("card not found")
		}
		return nil, err
	}

	return card, nil
}

// GetCardByCardNo fetches a card by its card number
func GetCardByCardNo(db *sql.DB, cardNo string) (*Card, error) {
	query := `SELECT id, card_no FROM cards WHERE card_no = ?`
	row := db.QueryRow(query, cardNo)

	var card Card
	err := row.Scan(&card.ID, &card.CardNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("card not found")
		}
		return nil, err
	}

	return &card, nil
}
