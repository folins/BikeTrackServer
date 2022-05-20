package service

import (
	"github.com/folins/biketrackserver"
	"github.com/folins/biketrackserver/package/repository"
)

type Authorization interface {
	CreateUser(user biketrackserver.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
}

type SMTP interface {
	SendConfirmCode(user biketrackserver.User) error
}

type BikeTrip interface {
	Create(userId int, trip biketrackserver.BikeTrip) (int, error)
	GetAll(userId int) ([]biketrackserver.BikeTrip, error)
	GetById(userId, tripId int) (biketrackserver.BikeTrip, error)
	Delete(userId, tripId int) error
}

type Service struct {
	Authorization
	BikeTrip
	SMTP
}

func NewService(repos *repository.Repository, smtp *SMTPService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		BikeTrip: NewBikeTripService(repos.BikeTrip),
		SMTP: smtp,
	}
}
