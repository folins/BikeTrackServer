package repository

import (
	"fmt"

	"github.com/folins/biketrackserver"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPosgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user biketrackserver.User) (int, error) {
	var id int 
	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
func (r *AuthPostgres) CreateTempUser(email string, code int) (int, error) {
	var id int 
	query := fmt.Sprintf("INSERT INTO %s (email, confirm_code) values ($1, $2) RETURNING id", tempUsersTable)
	row := r.db.QueryRow(query, email, code)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (biketrackserver.User, error) {
	var user biketrackserver.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}

func (r *AuthPostgres) GetTempUser(email, password string) (biketrackserver.tempUser, error) {
	var user biketrackserver.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}