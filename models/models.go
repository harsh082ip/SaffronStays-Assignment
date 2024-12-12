package models

import "time"

type Room struct {
	RoomID         int         `json:"room_id"`
	RatePerNight   float64     `json:"rate_per_night"`
	MaxGuests      int         `json:"max_guests"`
	AvailableDates []time.Time `json:"available_dates"` // Dates when the room is available
}

type NightRates struct {
	AverageRate float64 `json:"average_rate"` // Average rate for the next 30 days
	HighestRate float64 `json:"highest_rate"` // Highest rate for the next 30 days
	LowestRate  float64 `json:"lowest_rate"`  // Lowest rate for the next 30 days
}

type Occupancy struct {
	Month               string  `json:"month"`                // Month name, e.g., "December 2024"
	TotalDays           int     `json:"total_days"`           // Total days in the month within the next 5 months
	AvailableDays       int     `json:"available_days"`       // Available days in the month
	OccupancyPercentage float64 `json:"occupancy_percentage"` // Percentage of available days
}

type Metrics struct {
	Occupancy  []Occupancy `json:"occupancy"`   // Month-wise occupancy details
	NightRates NightRates  `json:"night_rates"` // Rates data for the next 30 days
}

type RoomResponse struct {
	RoomID  int     `json:"room_id"`
	Metrics Metrics `json:"metrics"`
}
