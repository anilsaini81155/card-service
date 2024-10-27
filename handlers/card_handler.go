// handlers/card.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anilsaini81155/card-service/models"
)

type CardStatusResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func CreateCard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserID   int    `json:"user_id"`
			MobileNo string `json:"mobile_no"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Fetch user by ID or Mobile No
		var user *models.User
		var err error
		if input.UserID != 0 {
			user, err = models.GetUserByID(db, input.UserID)
		} else if input.MobileNo != "" {
			user, err = models.GetUserByMobileNo(db, input.MobileNo)
		} else {
			http.Error(w, "Either user_id or mobile_no must be provided", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Generate a unique card number
		cardNo := models.GenerateCardNo()

		// Create the card for the user
		cardID, err := models.CreateCard(db, cardNo, user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"card_id": cardID,
			"card_no": cardNo,
			"message": "Card created successfully",
		}

		json.NewEncoder(w).Encode(response)
	}
}

func GetCardStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mobileNo := r.URL.Query().Get("mobile_no")
		cardID := r.URL.Query().Get("card_id")

		if mobileNo == "" && cardID == "" {
			// If both are empty, return an error
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CardStatusResponse{
				Error: "Invalid input: Either mobile_no or card_id must be provided.",
			})
			return
		}

		// var input struct {
		// 	CardID   int    `json:"card_id"`
		// 	MobileNo string `json:"mobile_no"`
		// }

		// if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// 	http.Error(w, "Invalid input", http.StatusBadRequest)
		// 	return
		// }

		var card *models.Card
		var err error

		// Fetch card using CardID or MobileNo
		if cardID != "" {
			card, err = models.GetCardByCardNo(db, cardID)
		} else if mobileNo != "" {
			user, err := models.GetUserByMobileNo(db, mobileNo)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			card, _ = models.GetCardByUserID(db, user.ID)
		} else {
			http.Error(w, "Either card_id or mobile_no must be provided", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Card not found", http.StatusNotFound)
			return
		}

		// Fetch delivery status by card ID
		delivery, err := models.GetDeliveryByCardID(db, card.ID)
		if err != nil {
			http.Error(w, "Delivery record not found", http.StatusNotFound)
			return
		}

		response := map[string]interface{}{
			"card_no": card.CardNo,
			"status":  delivery.Status,
			"message": "Status fetched successfully",
		}

		json.NewEncoder(w).Encode(response)
	}
}
