package biketrackserver

import (
	"time"
)

type BikeTrip struct {
	Id int    `json:"id" db:"id"`
    DateStart time.Time `json:"date_start" db:"date_start"`
    DateEnd time.Time `json:"date_end" db:"date_end"`
    PauseDuration int64 `json:"pause_duration" db:"pause_duration"`
    ActiveDuration int64 `json:"active_duration" db:"active_duration"`
    Distance float32 `json:"distance" db:"distance"`
    AvgSpeed float32 `json:"avg_speed" db:"avg_speed"`
    MaxSpeed float32 `json:"max_speed" db:"max_speed"`
    // UserId int `json:"-" db:"user_id"`
}

type TripPoint struct {
	Id int    `json:"id" db:"id"`
    Latitude float64 `json:"latitude" db:"latitude"`
    Longitude float64 `json:"longitude" db:"longitude"`
    Date time.Time `json:"date" db:"point_date"`
    Speed float32 `json:"speed" db:"speed"`
    // TripId int64 `json:"active_duration" db:"trip_id"`
}