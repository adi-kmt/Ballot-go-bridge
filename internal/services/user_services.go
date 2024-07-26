package services

import (
	"github.com/adi-kmt/voting-backend-go-grpc/internal/jwt"
	"github.com/adi-kmt/voting-backend-go-grpc/internal/messages"
	"github.com/adi-kmt/voting-backend-go-grpc/internal/repositories"
)

type UserService struct {
	repo repositories.IRepo
}

func NewUserService(repo repositories.IRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) AddUser(userName, password string, isStandingForElection bool) (string, *messages.AppError) {
	isUserPresent, err0 := s.repo.CheckUserExists(userName)
	if err0 != nil {
		return "", err0
	}
	if isUserPresent {
		return "", messages.NotFound("User already exists")
	}
	err := s.repo.AddUser(userName, password, isStandingForElection)
	if err != nil {
		return "", err
	}
	jwtToken, err1 := jwt.GenerateToken(userName)
	if err1 != nil {
		return "", err1
	}
	return jwtToken, nil
}

func (s *UserService) ValidateUser(userName, password string) (string, *messages.AppError) {
	// First checking if user is present, otherwise sending a 404 user not found
	isUserPresent, err0 := s.repo.CheckUserExists(userName)
	if err0 != nil {
		return "", err0
	}
	if !isUserPresent {
		return "", messages.NotFound("User not found")
	}
	// Now checking if the password is correct for found user
	userPassword, err1 := s.repo.GetUserPassword(userName)
	if err1 != nil {
		return "", err1
	}
	if userPassword != password {
		return "", messages.Unauthorized("Invalid password")
	}
	jwtToken, err2 := jwt.GenerateToken(userName)
	if err2 != nil {
		return "", err2
	}
	return jwtToken, nil
}
