package watcher

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/anilsaini81155/card-service/handlers"
	"github.com/fsnotify/fsnotify"
)

// WatchDirectory function to monitor the data directory
func WatchDirectory(dataDir string, db *sql.DB) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add directory to watch
	err = watcher.Add(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watching directory: ./data")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Check for creation of a new file
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Println("New file detected:", event.Name)

				// Check the file name and call the corresponding handler function
				switch filepath.Base(event.Name) {
				case "Pickup.csv":
					log.Println("Processing Pickup.csv")
					handlers.ProcessPickupFile(db) // Call the function to handle Pickup.csv
				case "Delivery_exceptions.csv":
					log.Println("Processing Delivery_exceptions.csv")
					handlers.ProcessDeliveryExceptionsFile(db) // Handle Delivery_exceptions.csv
				case "Delivered.csv":
					log.Println("Processing Delivered.csv")
					handlers.ProcessDeliveredFile(db) // Handle Delivered.csv
				case "Returned.csv":
					log.Println("Processing Returned.csv")
					handlers.ProcessReturnedFile(db) // Handle Returned.csv
				default:
					log.Println("Unknown file detected, skipping:", event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
