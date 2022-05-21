package service

import (
	"github.com/folins/biketrackserver"
	"github.com/folins/biketrackserver/package/repository"
)

type User interface {
	Create(user biketrackserver.User) (int, error)
	GetIdByEmail(email string) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
	Update(userId int, input biketrackserver.UserUpdateInput) error
	CheckPassword(userId int, password string) error
	CheckConfirmCode(email string, code int) error
	SetPassword(email, password string, code int) error
}

type SMTP interface {
	SendConfirmCode(user biketrackserver.User) error
	SendResetCode(email string, code int) error
}

type BikeTrip interface {
	Create(userId int, trip biketrackserver.BikeTrip) (int, error)
	GetAll(userId int) ([]biketrackserver.BikeTrip, error)
	GetById(userId, tripId int) (biketrackserver.BikeTrip, error)
	Delete(userId, tripId int) error
}

type TripPoint interface {
	Create(userId, tripId int, point biketrackserver.TripPoint) (int, error)
	GetAll(userId, tripId int) ([]biketrackserver.TripPoint, error)
	GetById(userId, tripId, pointId int) (biketrackserver.TripPoint, error)
	// Delete(userId, tripPointId int) error
}

type Service struct {
	User
	BikeTrip
	TripPoint
	SMTP
}

func NewService(repos *repository.Repository, smtp *SMTPService) *Service {
	return &Service{
		User:      NewUserService(repos.User),
		BikeTrip:  NewBikeTripService(repos.BikeTrip),
		TripPoint: NewTripPointService(repos.TripPoint, repos.BikeTrip),
		SMTP:      smtp,
	}
}
