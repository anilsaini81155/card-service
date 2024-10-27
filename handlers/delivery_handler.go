// handlers/delivery.go
package handlers

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/anilsaini81155/card-service/models"
)

func ProcessPickupFile(db *sql.DB) {
	file, err := os.Open("data/Pickup.csv")
	if err != nil {
		log.Fatal("Could not open the file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read() // Skip the header row
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading record:", err)
			continue
		}

		// Extract values from the CSV record
		cardNo := record[1]
		mobileNo := record[2]

		// Fetch the card and user from the database
		card, err := models.GetCardByCardNo(db, cardNo)
		if err != nil {
			log.Println("Card not found:", err)
			continue
		}

		user, err := models.GetUserByMobileNo(db, mobileNo)

		if err != nil {
			log.Println("User not found:", err)
			continue
		}
		log.Println("user found", user)

		// Create a delivery record with status "pickedup"
		_, err = models.CreateDelivery(db, card.ID, "pickedup")
		if err != nil {
			log.Println("Failed to create delivery record:", err)
			continue
		}
	}
}

func ProcessDeliveredFile(db *sql.DB) {
	file, err := os.Open("data/Delivered.csv")
	if err != nil {
		log.Fatal("Could not open the file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read() // Skip the header row
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading record:", err)
			continue
		}

		// Extract values from the CSV record
		cardNo := record[1]
		comment := record[4]

		// Fetch the card from the database
		card, err := models.GetCardByCardNo(db, cardNo)
		if err != nil {
			log.Println("Card not found:", err)
			continue
		}

		// Fetch delivery record and update status to "delivered"
		delivery, err := models.GetDeliveryByCardID(db, card.ID)
		if err != nil {
			log.Println("Delivery record not found:", err)
			continue
		}

		// Update the delivery status to "delivered"
		err = models.UpdateDeliveryStatus(db, delivery.ID, "delivered")
		if err != nil {
			log.Println("Failed to update delivery status:", err)
			continue
		}

		// Log the comment in the delivery_logs table
		err = models.CreateDeliveryLog(db, delivery.ID, comment)
		if err != nil {
			log.Println("Failed to create delivery log:", err)
			continue
		}
	}
}

// ProcessDeliveryExceptionsFile handles Delivery_exceptions.csv
func ProcessDeliveryExceptionsFile(db *sql.DB) {
	file, err := os.Open("./data/Delivery_exceptions.csv")
	if err != nil {
		log.Println("Error opening Delivery_exceptions.csv:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading Delivery_exceptions.csv:", err)
		return
	}

	for _, record := range records[1:] {
		cardNo := record[1]
		comment := record[4]

		card, cardErr := models.GetCardByCardNo(db, cardNo)

		if cardErr != nil {
			log.Println("Error fetching delivery for card:", cardNo, cardErr)
			continue
		}

		// Convert cardID from string to int
		// cardID, err := strconv.Atoi(cardIDStr)
		// if err != nil {
		// 	// Handle conversion error
		// 	log.Printf("Error converting cardID to int: %v", err)
		// 	continue
		// }

		// Get delivery record using card ID
		delivery, err := models.GetDeliveryByCardID(db, card.ID)
		if err != nil {
			log.Println("Error fetching delivery for card:", cardNo, err)
			continue
		}

		// Insert into delivery_logs table
		err = models.CreateDeliveryLog(db, delivery.ID, comment)
		if err != nil {
			log.Println("Error creating delivery log:", err)
			continue
		}

		// Update the delivery status
		var newStatus string
		if delivery.Status == "delivery_attempt_1" {
			newStatus = "delivery_attempt_2"
		} else if delivery.Status == "delivery_attempt_2" {
			newStatus = "Returned"
		} else {
			newStatus = "delivery_attempt_1"
		}

		err = models.UpdateDeliveryStatus(db, delivery.ID, newStatus)
		if err != nil {
			log.Println("Error updating delivery status:", err)
		}
	}
	log.Println("Finished processing Delivery_exceptions.csv")
}

// ProcessReturnedFile handles Returned.csv
func ProcessReturnedFile(db *sql.DB) {
	file, err := os.Open("./data/Returned.csv")
	if err != nil {
		log.Println("Error opening Returned.csv:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading Returned.csv:", err)
		return
	}

	for _, record := range records[1:] {
		fmt.Println("Records===>", record[1], record[2], record[3])
		cardNo := record[1]

		card, cardErr := models.GetCardByCardNo(db, cardNo)

		if cardErr != nil {
			log.Println("Error fetching Card :", cardNo, cardErr)
			continue
		}

		// Convert cardID from string to int
		// cardID, err := strconv.Atoi(cardIDStr)
		// if err != nil {
		// 	// Handle conversion error
		// 	log.Printf("Error converting cardID to int: %v", err)
		// 	continue
		// }

		// Get delivery record using card ID
		delivery, err := models.GetDeliveryByCardID(db, card.ID)
		if err != nil {
			log.Println("Error fetching delivery for card:", cardNo, err)
			continue
		}

		// Update delivery status to Returned
		err = models.UpdateDeliveryStatus(db, delivery.ID, "Returned")
		if err != nil {
			log.Println("Error updating delivery status to Returned:", err)
			continue
		}
	}
	log.Println("Finished processing Returned.csv")
}
