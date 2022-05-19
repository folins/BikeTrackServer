package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)
	
const (
	usersTable = "users"
	tempUsersTable = "temp_users"
)


type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBNmae   string
	SSLMode  string
}

func NewPostgreDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	cfg.Host, cfg.Port, cfg.Username, cfg.DBNmae, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}