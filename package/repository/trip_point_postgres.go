package repository

import (
	"fmt"

	"github.com/folins/biketrackserver"
	"github.com/jmoiron/sqlx"
)

type TripPointPostgres struct {
	db *sqlx.DB
}

func NewTripPointPostgres(db *sqlx.DB) *TripPointPostgres {
	return &TripPointPostgres{db: db}
}

func (r *TripPointPostgres) Create(tripId int, point biketrackserver.TripPoint) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (latitude, longitude, 
						point_date, speed, trip_id)
						VALUES ($1, $2, $3, $4, $5) RETURNING id`, tripPointsTable)
	row := r.db.QueryRow(query, point.Latitude, point.Longitude, point.Date, point.Speed, tripId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TripPointPostgres) GetAll(tripId int) ([]biketrackserver.TripPoint, error) {
	var trips []biketrackserver.TripPoint

	query := fmt.Sprintf(`SELECT id, latitude, longitude, point_date, speed
						FROM %s WHERE trip_id = $1`, tripPointsTable)
	err := r.db.Select(&trips, query, tripId)

	return trips, err
}

func (r *TripPointPostgres) GetById(tripId, pointId int) (biketrackserver.TripPoint, error) {
	var trip biketrackserver.TripPoint

	query := fmt.Sprintf(`SELECT id, latitude, longitude, point_date, speed
						FROM %s WHERE trip_id = $1 AND id = $2`, tripPointsTable)
	err := r.db.Get(&trip, query, tripId, pointId)

	return trip, err
}
