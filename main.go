package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/anilsaini81155/card-service/db"
	"github.com/anilsaini81155/card-service/handlers"
	"github.com/anilsaini81155/card-service/watcher"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Get the singleton DB instance
	database := db.GetDB()
	// Create router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/users", handlers.CreateUser(database)).Methods("POST")
	r.HandleFunc("/cards", handlers.CreateCard(database)).Methods("POST")
	r.HandleFunc("/card-status", handlers.GetCardStatus(database)).Methods("GET")

	// Watch the 'data' directory for file changes
	dataDir := "./data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		log.Fatalf("Data directory %s does not exist", dataDir)
	}

	// Use WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		watcher.WatchDirectory(dataDir, database)
	}()

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

	// Wait for fsnotify goroutine to finish
	wg.Wait()
}
