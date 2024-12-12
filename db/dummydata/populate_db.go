package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/lib/pq"
)

const maxBatchSize = 2175 // Max rows per batch based on PostgreSQL's 65535 parameter limit

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

	var bulkInsertValues []string
	var bulkInsertArgs []interface{}

	for i := 0; i < 1500; i++ {
		availableDates, ratePerNight := generateDatesAndRates()

		maxGuests := rand.Intn(6) + 1 // Max guests between 1 and 6

		for j := 0; j < len(availableDates); j++ {
			// Add each row of data to bulkInsertValues
			placeholder := fmt.Sprintf("($%d, $%d, $%d)", len(bulkInsertArgs)+1, len(bulkInsertArgs)+2, len(bulkInsertArgs)+3)
			bulkInsertValues = append(bulkInsertValues, placeholder)

			// Append values to the argument slice
			bulkInsertArgs = append(bulkInsertArgs, pq.Array(ratePerNight), maxGuests, pq.Array(availableDates))
		}

		if len(bulkInsertValues) >= maxBatchSize {
			err := executeBulkInsert(db, bulkInsertValues, bulkInsertArgs)
			if err != nil {
				log.Fatalf("Error performing bulk insert: %v", err)
			}

			// Reset for the next batch
			bulkInsertValues = nil
			bulkInsertArgs = nil
		}
	}

	if len(bulkInsertValues) > 0 {
		err := executeBulkInsert(db, bulkInsertValues, bulkInsertArgs)
		if err != nil {
			log.Fatalf("Error performing bulk insert: %v", err)
		}
	}

	fmt.Println("Dummy data inserted successfully!")
}

func executeBulkInsert(db *sql.DB, bulkInsertValues []string, bulkInsertArgs []interface{}) error {

	query := fmt.Sprintf("INSERT INTO rooms (rate_per_night, max_guests, available_dates) VALUES %s", join(bulkInsertValues, ","))

	fmt.Println("Executing bulk insert with values:", bulkInsertArgs)

	_, err := db.Exec(query, bulkInsertArgs...)
	return err
}

func generateDatesAndRates() ([]time.Time, []float64) {
	var dates []time.Time
	var rates []float64
	startDate := time.Now()
	endDate := startDate.AddDate(0, 5, 0)

	// Randomly generate up to 30 available dates
	numDates := rand.Intn(30) + 1
	for i := 0; i < numDates; i++ {
		randomDate := startDate.AddDate(0, 0, rand.Intn(int(endDate.Sub(startDate).Hours()/24)))
		dates = append(dates, randomDate)

		// Generate a corresponding nightly rate for each date
		randomRate := rand.Float64()*200 + 50 // Rates between 50 and 250

		// Round off randomRate to upto 2 decimal places
		randomRate = math.Round(randomRate*100) / 100
		rates = append(rates, randomRate)
	}

	return dates, rates
}

func join(values []string, separator string) string {
	result := ""
	for i, value := range values {
		if i > 0 {
			result += separator
		}
		result += value
	}
	return result
}
