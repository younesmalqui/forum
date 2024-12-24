package services

import (
	"errors"
	"strings"

	"forum/config"
	"forum/models"
	"forum/utils"
)

var (
	errInvalidUserPass  = errors.New("invalid username or password")
	errUserOrEmailExist = errors.New("username or email already used")
	errFieldsEmpty      = errors.New("all fields must be completed")
)

func RegisterUser(user *models.User) error {
	userRepo := models.NewUserRepository()

	// check if the username or email alread yexist
	isUserExist, err := userRepo.UserExists(user.Username, user.Email)
	if err != nil {
		return err
	}
	if isUserExist {
		return config.NewError(errUserOrEmailExist)
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return config.NewInternalError(err)
	}
	user.Password = hashedPassword
	return userRepo.CreateUser(user)
}

func LoginUser(username, password string) (*models.User, error) {
	if len(strings.TrimSpace(username)) == 0 || len(strings.TrimSpace(password)) == 0 {
		return nil, config.NewError(errFieldsEmpty)
	}
	userRepo := models.NewUserRepository()
	user, err := userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if utils.CheckPassword(user.Password, password) != nil {
		return nil, config.NewError(errInvalidUserPass)
	}
	return user, nil
}
