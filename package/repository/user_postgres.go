package repository

import (
	"fmt"
	"strings"

	"github.com/folins/biketrackserver"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPosgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(user biketrackserver.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash, confirm_code) values ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Email, user.Password, user.ConfirmCode)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) Get(email, password string) (biketrackserver.User, error) {
	var user biketrackserver.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}

func (r *UserPostgres) CheckConfirmCode(email string, code int) error {
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND confirm_code = $2", usersTable)
	_, err := r.db.Query(query, email, code)

	return err
}

func (r *UserPostgres) GetIdByEmail(email string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1", usersTable)
	err := r.db.Get(&id, query, email)

	return id, err

}
func (r *UserPostgres) GetIdByEmailAndConfirmCode(email string, code int) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND confirm_code = $2", usersTable)
	err := r.db.Get(&id, query, email, code)

	return id, err
}

func (r *UserPostgres) CheckPassword(userId int, password string) error {
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1 AND password_hash = $2", usersTable)
	_, err := r.db.Query(query, userId, password)

	return err
}

func (r *UserPostgres) Update(userId int, input biketrackserver.UserUpdateInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *input.Email)
		argId++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password_hash=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	if input.ConfirmCode != nil {
		setValues = append(setValues, fmt.Sprintf("confirm_code=$%d", argId))
		args = append(args, *input.ConfirmCode)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", usersTable, setQuery, argId)
	args = append(args, userId)

	logrus.Debugf("SetQuery: %s", setQuery)

	_, err := r.db.Exec(query, args...)
	return err
}
