package userservice

import (
	"errors"
	"orb-api/models"
	repository "orb-api/repositories"
	"orb-api/repositories/user"
)

func SetupUserService(repo repository.Repository) *UserService {
	return &UserService{
		UserRepo: repo.User,
	}
}

func (service *UserService) CreateNewUser(credentials ICreateUser) (*models.User, error) {
	// Check if the email is not being used by anyone else
	userArray, readErr := service.UserRepo.ReadBy(user.IReadBy{
		Email: &credentials.Email,
	})

	if readErr != nil {
		return nil, readErr
	}

	if len(userArray) == 1 {
		return nil, errors.New("This email is already being used")
	}

	// Check if the username is not being used by anyone else
	userArray, readErr = service.UserRepo.ReadBy(user.IReadBy{
		Name: &credentials.Name,
	})

	if readErr != nil {
		return nil, readErr
	}

	if len(userArray) == 1 {
		return nil, errors.New("This username is already being used")
	}

	// Check email, username and password length
	if !user.ValidUserName(credentials.Name) {
		return nil, errors.New("Invalid username size")
	}

	if !user.ValidUserEmail(credentials.Email) {
		return nil, errors.New("Invalid email size")
	}

	if !user.ValidUserPassword(credentials.Password) {
		return nil, errors.New("Invalid password size")
	}

	// Hash the password to prevent security vulnerabilities
	hashedPassword, hashErr := HashPassword(credentials.Password)

	if hashErr != nil {
		return nil, hashErr
	}

	newUser, createErr := service.UserRepo.Create(user.ICreate{
		Name:     credentials.Name,
		Email:    credentials.Email,
		Password: hashedPassword,
		Status:   models.UStatusProcessing,
	})

	if createErr != nil {
		return nil, createErr
	}

	return newUser, nil
}

func (service *UserService) UpdateEmail(id uint, email string) (*models.User, error) {
	// Check if the id belongs a valid user
	if !service.UserRepo.ValidUser(id) {
		return nil, errors.New("Invalid user id")
	}

	// Check if the email is not being used by anyone else and different by current
	userArray, readErr := service.UserRepo.ReadBy(user.IReadBy{
		Email: &email,
	})

	if readErr != nil {
		return nil, readErr
	}

	if len(userArray) > 0 {
		return nil, errors.New("This email is already being used")
	}

	// Update e-mail
	userUpdate, updateErr := service.UserRepo.Update(user.IUpdate{
		ID:    id,
		Email: &email,
	})

	if updateErr != nil {
		return nil, updateErr
	}

	return userUpdate, nil
}

func (service *UserService) UpdateStatus(id uint, status uint) (*models.User, error) {
	// Check if the status is valid
	if !user.ValidUserStatus(status) {
		return nil, errors.New("Invalid status")
	}

	// Check if the id belongs a valid user
	if !service.UserRepo.ValidUser(id) {
		return nil, errors.New("Invalid user id")
	}

	// Update status
	userUpdate, updateErr := service.UserRepo.Update(user.IUpdate{
		ID:     id,
		Status: &status,
	})

	if updateErr != nil {
		return nil, updateErr
	}

	return userUpdate, nil
}
