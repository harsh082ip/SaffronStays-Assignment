package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/lib/pq"
)

func main() {

	serviceURI := os.Getenv("SERVICE_URI")
	if serviceURI == "" {
		log.Fatal("SERVICE_URI is not set. Please provide the database connection string.")
	}

	db, err := sql.Open("postgres", serviceURI)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	for i := 0; i < 1500; i++ {
		ratePerNight := rand.Float64()*200 + 50 // Rates between 50 and 250
		maxGuests := rand.Intn(6) + 1           // Max guests between 1 and 6
		availableDates := generateRandomDates()

		query := `INSERT INTO rooms (rate_per_night, max_guests, available_dates) VALUES ($1, $2, $3)`
		_, err := db.Exec(query, ratePerNight, maxGuests, pq.Array(availableDates))
		if err != nil {
			log.Printf("Error inserting row %d: %v", i+1, err)
		}
	}

	fmt.Println("Dummy data inserted successfully!")
}

func generateRandomDates() []time.Time {
	var dates []time.Time
	startDate := time.Now()
	endDate := startDate.AddDate(0, 5, 0)

	// Randomly generate up to 30 available dates
	numDates := rand.Intn(30) + 1
	for i := 0; i < numDates; i++ {
		randomDate := startDate.AddDate(0, 0, rand.Intn(int(endDate.Sub(startDate).Hours()/24)))
		dates = append(dates, randomDate)
	}

	return dates
}
