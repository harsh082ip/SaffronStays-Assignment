package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/SaffronStays-Assignment/db"
	"github.com/harsh082ip/SaffronStays-Assignment/models"
	"github.com/lib/pq"
)

func MetricsControllers(c *gin.Context) {

	room_id := c.Param("room_id")
	room_id_int, err := RoomStrToInt(room_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "str to int conversion failed:" + err.Error(),
		})
		return
	}
	if room_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"staus": http.StatusBadRequest,
			"error": "/room_id cannot be empty",
		})
		return
	}
	var room models.Room
	db, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "error connecting to postgres:" + err.Error(),
		})
		return
	}

	room, err = GetRoomByID(db, room_id_int)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error in getting data from DB:" + err.Error(),
		})
		return
	}

	metrics := CalculateMetrics(room)

	response := models.RoomResponse{
		RoomID:  room.RoomID,
		Metrics: metrics,
	}

	c.JSON(http.StatusOK, response)

}

func CalculateMetrics(room models.Room) models.Metrics {
	// Nightrates
	var totalRate float64
	var highestRate float64 = room.RatePerNight[0]
	var lowestRate float64 = room.RatePerNight[0]

	for _, rate := range room.RatePerNight {
		totalRate += rate
		if rate > highestRate {
			highestRate = rate
		}
		if rate < lowestRate {
			lowestRate = rate
		}
	}

	averageRate := totalRate / float64(len(room.RatePerNight))

	nightRates := models.NightRates{
		AverageRate: averageRate,
		HighestRate: highestRate,
		LowestRate:  lowestRate,
	}

	// Occupancy
	occupancyMap := make(map[string]*models.Occupancy)

	for _, date := range room.AvailableDates {
		month := date.Format("January 2006")
		if _, exists := occupancyMap[month]; !exists {
			occupancyMap[month] = &models.Occupancy{
				Month: month,
			}
		}
		occupancyMap[month].AvailableDays++
	}

	occupancy := make([]models.Occupancy, 0, len(occupancyMap))
	for _, occ := range occupancyMap {
		occ.TotalDays = DaysInMonth(occ.Month)
		// occ.OccupancyPercentage = (float64(occ.AvailableDays) / float64(occ.TotalDays)) * 100
		occ.OccupancyPercentage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(occ.AvailableDays)/float64(occ.TotalDays))*100), 64)
		occupancy = append(occupancy, *occ)
	}

	return models.Metrics{
		NightRates: nightRates,
		Occupancy:  occupancy,
	}
}

func DaysInMonth(month string) int {
	t, _ := time.Parse("January 2006", month)
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func GetRoomByID(db *sql.DB, roomID int) (models.Room, error) {
	var room models.Room

	query := "SELECT room_id, rate_per_night, max_guests, available_dates FROM rooms WHERE room_id = $1"

	// Execute the query and scan into the room struct
	err := db.QueryRow(query, roomID).Scan(
		&room.RoomID,
		pq.Array(&room.RatePerNight),
		&room.MaxGuests,
		&room.AvailableDatesRaw, // Temporarily store raw date strings
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return room, fmt.Errorf("no room found with room_id %d", roomID)
		}
		return room, fmt.Errorf("error retrieving room: %v", err)
	}

	log.Print("dateStr")
	// Convert AvailableDatesRaw (strings) to time.Time and append them
	for _, dateStr := range room.AvailableDatesRaw {
		// Parse the raw date string to time.Time
		log.Print(dateStr)
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return models.Room{}, fmt.Errorf("error parsing date %s: %v", dateStr, err)
		}
		dateStr = date.Format("2006-01-02")
		finalDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return models.Room{}, fmt.Errorf("error parsing date %s: %v", dateStr, err)
		}
		log.Print("finalStr", finalDate)
		// Append the parsed time.Time object to AvailableDates
		room.AvailableDates = append(room.AvailableDates, date)
	}

	log.Print("dateStr2")
	return room, nil
}

func RoomStrToInt(roomID string) (int, error) {
	roomInt, err := strconv.Atoi(roomID)
	if err != nil {
		return 0, err
	}
	return roomInt, err
}
