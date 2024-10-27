package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	MobileNo  string    `json:"mobile_no"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUser creates a new user in the database.
func CreateUser(db *sql.DB, mobileNo string) (int64, error) {
	// Validate for 9-digit mobile number
	if len(mobileNo) != 9 {
		return 0, errors.New("mobile number must be 9 digits")
	}

	// Check for duplicate mobile number
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE mobile_no = ?)", mobileNo).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("mobile number already exists")
	}

	// Insert the user into the database
	result, err := db.Exec("INSERT INTO users (mobile_no, created_at, updated_at) VALUES (?, ?, ?)", mobileNo, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetUserByID fetches a user by their ID.
func GetUserByID(db *sql.DB, userID int) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT id, mobile_no FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.MobileNo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByMobileNo fetches a user by their mobile number.
func GetUserByMobileNo(db *sql.DB, mobileNo string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT id, mobile_no FROM users WHERE mobile_no = ?", mobileNo).Scan(
		&user.ID, &user.MobileNo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
