package services

import (
	"github.com/adi-kmt/ai-streak-backend-go/internal/messages"
	"github.com/adi-kmt/ai-streak-backend-go/internal/repositories"
)

type UserService struct {
	repo repositories.IRepo
}

func NewUserService(repo repositories.IRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) AddUser(userName, password string, isStandingForElection bool) (string, error) {
	err := s.repo.AddUser(userName, password, isStandingForElection)
	if err != nil {
		return "", err
	}
	// We will send the token in the response
	return "", nil
}

func (s *UserService) ValidateUser(userName, password string) (string, error) {
	isUserPresent, err0 := s.repo.CheckUserExists(userName)
	if err0 != nil {
		return "", err0
	}
	if !isUserPresent {
		return "", messages.NotFound("User not found")
	}
	userPassword, err1 := s.repo.GetUserPassword(userName)
	if err1 != nil {
		return "", err1
	}
	if userPassword != password {
		return "", messages.Unauthorized("Invalid password")
	}
	// We will send the token in the response
	return "", nil
}
