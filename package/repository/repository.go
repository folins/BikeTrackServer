package repository

import (
	"github.com/folins/biketrackserver"
	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(user biketrackserver.User) (int, error)
	Get(email, password string) (biketrackserver.User, error)
	GetIdByEmail(email string) (int, error)
	GetUserByEmail(email string) (biketrackserver.User, error)
	GetIdByEmailAndConfirmCode(email string, code int) (int, error)
	Update(userId int, input biketrackserver.UserUpdateInput) error
	CheckPassword(userId int, password string) error
	CheckConfirmCode(email string, code int) error
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
	User
	BikeTrip
	TripPoint
}


func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserPosgres(db),
		BikeTrip: NewBikeTripPostgres(db),
		TripPoint: NewTripPointPostgres(db),
	}
}