package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/folins/biketrackserver"
	"github.com/folins/biketrackserver/package/repository"
)

const (
	salt      = "sadfashdgfbjkh435asdfa3"
	signinKey = "pwoeijfgsdDf#$4asdf5%"
	tokenTTL  = 24 * 7 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type UserService struct {
	repos repository.User
}

func NewUserService(repos repository.User) *UserService {
	return &UserService{repos: repos}
}

func (s *UserService) Create(email string, code int) (int, error) {

	user, err := s.repos.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			var inputUser biketrackserver.User
			strCode := strconv.Itoa(code)
			inputUser.Password = generatePasswordHash(strCode)
			inputUser.ConfirmCode = code
			inputUser.Email = email
			inputUser.IsRegistered = false
			return s.repos.Create(inputUser)
		}
	}

	if user.IsRegistered == false {
		var inputUser biketrackserver.UserUpdateInput
		strCode := strconv.Itoa(code)
		passwordHash := generatePasswordHash(strCode)
		*inputUser.Password = passwordHash
		*inputUser.ConfirmCode = code
		*inputUser.Email = email
		return 0, s.repos.Update(user.Id, inputUser)
	}

	return 0, errors.New("Server error")
}

func (s *UserService) GetIdByEmail(email string) (int, error) {
	return s.repos.GetIdByEmail(email)
}

func (s *UserService) GenerateToken(email, password string) (string, error) {
	user, err := s.repos.Get(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signinKey))
}

func (s *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signinKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *UserService) Update(userId int, input biketrackserver.UserUpdateInput) error {
	if input.Password != nil {
		*input.Password = generatePasswordHash(*input.Password)
	}
	return s.repos.Update(userId, input)
}

func (s *UserService) CheckPassword(userId int, password string) error {
	return s.repos.CheckPassword(userId, generatePasswordHash(password))
}

func (s *UserService) CheckConfirmCode(email string, code int) error {
	id, err := s.repos.GetIdByEmailAndConfirmCode(email, code)
	if err == nil {
		var user biketrackserver.UserUpdateInput
		*user.IsRegistered = true
		s.repos.Update(id, user)
	}
	return err
}

func (s *UserService) CheckEmailExistence(email string) (bool, error) {
	_, err := s.repos.GetIdByEmail(email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *UserService) SetPassword(email, password string, code int) error {
	id, err := s.repos.GetIdByEmailAndConfirmCode(email, code)
	if err != nil {
		return err
	}

	var input biketrackserver.UserUpdateInput
	passwordHash := generatePasswordHash(password)
	input.Password = &passwordHash
	*input.ConfirmCode = 0
	*input.IsRegistered = true

	return s.repos.Update(id, input)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
