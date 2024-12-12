package models

import "time"

type Room struct {
	RoomID         int         `json:"room_id"`
	RatePerNight   float64     `json:"rate_per_night"`
	MaxGuests      int         `json:"max_guests"`
	AvailableDates []time.Time `json:"available_dates"`
}

type NightRates struct {
	AverageRate float64 `json:"average_rate"`
	HighestRate float64 `json:"highest_rate"`
	LowestRate  float64 `json:"lowest_rate"`
}

type Metrics struct {
	OccupancyPercentage float64 `json:"occupancy_percentage"`
	NightRates          NightRates
}

type RoomResponse struct {
	RoomID  int `json:"room_id"`
	Metrics Metrics
}
