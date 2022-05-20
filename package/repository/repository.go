package repository

import (
	"github.com/folins/biketrackserver"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user biketrackserver.User) (int, error)
	GetUser(email, password string) (biketrackserver.User, error)
}

type BikeTrip interface {
	Create(userId int, trip biketrackserver.BikeTrip) (int, error)
	GetAll(userId int) ([]biketrackserver.BikeTrip, error)
	GetById(userId, tripId int) (biketrackserver.BikeTrip, error)
	Delete(userId, tripId int) error
}

type TripPoint interface {
	Create(tripId int, trip biketrackserver.TripPoint) (int, error)
	GetAll(tripId int) ([]biketrackserver.TripPoint, error)
	GetById(tripId, pointId int) (biketrackserver.TripPoint, error)
}

type Repository struct {
	Authorization
	BikeTrip
	TripPoint
}


func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPosgres(db),
		BikeTrip: NewBikeTripPostgres(db),
		TripPoint: NewTripPointPostgres(db),
	}
}