package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	once     sync.Once
	instance *sql.DB
)

func GetDB() *sql.DB {
	once.Do(func() {
		dsn := "user:password@tcp(127.0.0.1:3306)/card_service"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening DB: %v", err)
		}

		// Check if the database connection is alive
		if err := db.Ping(); err != nil {
			log.Fatalf("Error pinging DB: %v", err)
		}

		instance = db
		fmt.Println("DB connection established")
	})

	return instance
}
