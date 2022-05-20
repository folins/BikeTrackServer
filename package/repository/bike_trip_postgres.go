package repository

import (
	"fmt"

	"github.com/folins/biketrackserver"
	"github.com/jmoiron/sqlx"
)

type BikeTripPostgres struct {
	db *sqlx.DB
}

func NewBikeTripPostgres(db *sqlx.DB) *BikeTripPostgres {
	return &BikeTripPostgres{db: db}
}

func (r *BikeTripPostgres) Create(userId int, trip biketrackserver.BikeTrip) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (date_start, date_end, pause_duration, 
						active_duration, distance, avg_speed, max_speed, user_id)
						 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, bikeTripsTable)
	row := r.db.QueryRow(query, trip.DateStart, trip.DateEnd, trip.PauseDuration,
						 trip.ActiveDuration, trip.Distance, trip.AvgSpeed, trip.MaxSpeed, userId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *BikeTripPostgres) GetAll(userId int) ([]biketrackserver.BikeTrip, error) {
	var trips []biketrackserver.BikeTrip

	query := fmt.Sprintf(`SELECT id, date_start, date_end, 
						pause_duration, active_duration, distance, avg_speed, max_speed 
						FROM %s WHERE user_id = $1`, bikeTripsTable)
	err := r.db.Select(&trips, query, userId)

	return trips, err
}

func (r *BikeTripPostgres) GetById(userId, tripId int) (biketrackserver.BikeTrip, error) {
	var trip biketrackserver.BikeTrip

	query := fmt.Sprintf(`SELECT id, date_start, date_end, pause_duration, active_duration, 
						distance, avg_speed, max_speed FROM %s 
						WHERE user_id = $1 AND id = $2`, bikeTripsTable)
	err := r.db.Get(&trip, query, userId, tripId)

	return trip, err
}

func (r *BikeTripPostgres) Delete(userId, tripId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", usersTable)
	_, err := r.db.Exec(query, tripId, userId)

	return err
}