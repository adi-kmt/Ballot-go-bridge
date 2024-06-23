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
	hasUserVoted, err2 := s.repo.CheckUserHasVoted(userName)
	if err2 != nil {
		return err2
	}
	if hasUserVoted {
		return messages.BadRequest("User has already voted")
	}
	// check if user is standing for election
	isUserStandingForElection, err1 := s.repo.CheckUserIsStandingForElection(VotedForUserName)
	if err1 != nil {
		return err1
	}
	if !isUserStandingForElection {
		return messages.BadRequest(fmt.Sprintf("User %s is not standing for election", VotedForUserName))
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
