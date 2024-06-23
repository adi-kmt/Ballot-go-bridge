package services

import (
	"fmt"

	"github.com/adi-kmt/ai-streak-backend-go/internal/entities"
	"github.com/adi-kmt/ai-streak-backend-go/internal/messages"
	"github.com/adi-kmt/ai-streak-backend-go/internal/repositories"
)

type VotingService struct {
	repo repositories.IRepo
}

func NewVotingService(repo repositories.IRepo) *VotingService {
	return &VotingService{
		repo: repo,
	}
}

func (s *VotingService) AddVote(userName, VotedForUserName string) *messages.AppError {
	// check if user exists
	isUserPresent, err0 := s.repo.CheckUserExists(userName)
	if err0 != nil {
		return err0
	}
	if !isUserPresent {
		return messages.NotFound(fmt.Sprintf("User %s does not exist", userName))
	}
	// check if user is standing for election
	isUserStandingForElection, err1 := s.repo.CheckUserIsStandingForElection(userName)
	if err1 != nil {
		return err1
	}
	if !isUserStandingForElection {
		return messages.BadRequest(fmt.Sprintf("User %s is not standing for election", userName))
	}
	err := s.repo.AddVote(userName, VotedForUserName)
	if err != nil {
		return err
	}
	return nil
}

func (s *VotingService) GetCurrentVoteSapshot() ([]entities.LeaderBoardItem, *messages.AppError) {
	return s.repo.GetCurrentVoteSapshot()
}
