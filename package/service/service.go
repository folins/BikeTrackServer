package service

import (
	"github.com/folins/biketrackserver"
	"github.com/folins/biketrackserver/package/repository"
)

type Authorization interface {
	CreateUser(user biketrackserver.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
	CreateTempUser(email string, code int) (int, error)
	VerifyConfirmCode(email string, code int) (bool, error)
}

type SMTP interface {
	SendConfirmCode(user biketrackserver.User, code int) error
}

type Service struct {
	Authorization
	SMTP
}

func NewService(repos *repository.Repository, smtp *SMTPService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		SMTP: smtp,
	}
}
