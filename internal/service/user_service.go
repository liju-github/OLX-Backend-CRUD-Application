package service

import (
	"context"
	"errors"

	"github.com/liju-github/internal/model"
	"github.com/liju-github/internal/repository"
	"github.com/liju-github/internal/utils"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func (service *UserService) RegisterUser(user model.User) error {
	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Check if the user already exists
	existingUser, err := service.UserRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user already exists") // User exists
	}
	// Add the user to the repository
	return service.UserRepo.AddUser(user)
}

func (service *UserService) Login(email, password string) (*model.User, error) {
	// Retrieve user by email
	user, err := service.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password") // Handle email not found
	}

	// Check if the password is correct
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid email or password") // Handle wrong password
	}

	return user, nil
}

func (service *UserService) AllUsers() ([]model.User, error) {
	users, err := service.UserRepo.GetAllUsers()
	if err != nil {
		return nil, errors.New(err.Error()) // Handle unexpected errors from repository
	}
	return users, nil
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.UserRepo.GetUserByEmail(email)
}

func (service *UserService) UpdateUserImage(ctx context.Context, userEmail string, newImageUrl string) error {
	return service.UserRepo.UpdateUserImage(ctx, userEmail, newImageUrl)
}
