// handlers/user.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anilsaini81155/card-service/models"
)

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			MobileNo string `json:"mobile_no"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		userID, err := models.CreateUser(db, input.MobileNo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := map[string]interface{}{
			"user_id": userID,
			"message": "User created successfully",
		}

		json.NewEncoder(w).Encode(response)
	}
}
