package repository

import (
	"github.com/folins/biketrackserver"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user biketrackserver.User) (int, error)
	GetUser(email, password string) (biketrackserver.User, error)
	CreateTempUser(email string, code int) (int, error)
	GetTempUser(email string, code int) (int, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPosgres(db),
	}
}